import React, { useEffect, useState } from 'react';
import { Switch, Route, useRouteMatch, useHistory, useParams, Link } from "react-router-dom";
import { useDispatch, useSelector } from 'react-redux';
import { selectItems, fetchItems, createItem, deleteItem } from '../store/itemsSlice';
import { Item } from '../model/item';
import { Report } from '../components/Report';
import axios from 'axios';

export function Items() {
    let { path } = useRouteMatch();
    return (
        <div>
            <Switch>
                <Route exact path={`${path}`}>
                    <ItemsIndex />
                </Route>
                <Route path={`${path}/create`}>
                    <AddItem />
                </Route>
                <Route path={`${path}/qr/:id`}>
                    <QR />
                </Route>
                <Route path={`${path}/report`}>
                    <Report />
                </Route>
            </Switch>
        </div>
    );
}

interface AddItemState {
    name: string,
    desired: string,
    actual: string
}

function AddItem() {
    const dispatch = useDispatch();
    const history = useHistory();

    const [values, setValues] = useState<AddItemState>({
        name: "",
        desired: "",
        actual: ""
    });

    const handleChange = (fieldName: keyof AddItemState) => (
        e: React.ChangeEvent<HTMLInputElement>
    ) => {
        setValues({ ...values, [fieldName]: e.currentTarget.value });
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        var item: Item = {
            name: values.name,
            desired: Number(values.desired),
            actual: Number(values.actual)
        }

        dispatch(createItem(item));
        history.push("/items");
    };

    return (
        <div className="row mt-4">
            <div className="col-4">
                <div className="card">
                    <h5 className="card-header">Create item</h5>
                    <div className="card-body">
                        <form>
                            <div className="mb-3">
                                <label className="form-label">Name</label>
                                <input className="form-control" name="name" value={values.name} onChange={handleChange("name")} />
                            </div>
                            <div className="mb-3">
                                <label className="form-label">Desired</label>
                                <input className="form-control" name="desired" value={values.desired} onChange={handleChange("desired")} />
                            </div>
                            <div className="mb-3">
                                <label className="form-label">Actual</label>
                                <input className="form-control" name="actual" value={values.actual} onChange={handleChange("actual")} />
                            </div>
                            <button className="btn btn-outline-primary" onClick={handleSubmit}>Add</button>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    );
}

interface ItemProps {
    item: Item
}

export function ItemRow(props: ItemProps) {
    const item = props.item
    const dispatch = useDispatch();

    return (
        <tr>
            <td>{item.id}</td>
            <td>{item.name}</td>
            <td>{item.desired}</td>
            <td>{item.actual}</td>
            <td>
                <div className="float-end">
                    <Link className="btn btn-link" to={`/items/qr/${item.id}`}><i className="fas fa-qrcode"></i></Link>
                    <button className="btn btn-link" onClick={() => dispatch(deleteItem(item))}><i className="fas fa-trash"></i></button>
                </div>
            </td>
        </tr>
    );
}

export function ItemsIndex() {
    const items = useSelector(selectItems);
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch(fetchItems());
    }, [dispatch]);

    return (
        <div>
            <div className="row mt-4">
                <div className="col-12">
                    <table className="table">
                        <thead>
                            <tr>
                                <th>#</th>
                                <th>Name</th>
                                <th>Desired</th>
                                <th>Actual</th>
                                <th></th>
                            </tr>
                        </thead>
                        <tbody className="align-middle">
                            {items.items.map((item: Item) => (
                                <ItemRow key={item.id} item={item} />
                            ))}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    );
}

interface QRParams {
    id: string
}

export function QR() {
    const { id } = useParams<QRParams>();
    const [item, setItem] = useState<Item>({
        name: "Loading",
        actual: 0,
        desired: 0
    });

    useEffect(() => {
        axios.get("http://localhost:8080/api/items/" + id).then(result => {
            setItem(result.data);
        }).catch(error => {
            console.log(error);
        });
    }, [id]);

    var QRCode = require('qrcode.react');

    return (
        <div className="row">
            <div className="col-3">
                <div className="card mt-4">
                    <div className="card-body">
                        <h4 className="card-title">{ item.name }</h4>
                        <div className="d-flex justify-content-center">
                            <QRCode className="d-flex justify-content-center" size={256} value={`http://stoqr.veicot.io/api/items/mark/${item.id}`} />
                        </div>
                        <p className="float-end mt-2"><i>generated by STOQR</i></p>                        
                    </div>
                </div>
            </div>
        </div>
    );
}