import React, { useEffect, useState } from 'react';
import { Switch, Route, useRouteMatch, useHistory, useParams, Link } from "react-router-dom";
import { useDispatch, useSelector } from 'react-redux';
import { selectItems, fetchItems, createItem, deleteItem, withdraw } from '../store/itemsSlice';
import { Item } from '../model/item';
import { axiosInstance } from '../service/service';
import { AxiosResponse, AxiosError } from 'axios';

export function Items() {
    let { path } = useRouteMatch();
    const dispatch = useDispatch();

    const onSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
        dispatch(fetchItems(e.currentTarget.value));
    };

    return (
        <div>
            <div className="row mt-4">
                <div className="col-9">
                    <h2>Items</h2>
                </div>
                <div className="col">
                    <div className="d-flex">
                        <div className="flex-grow-1">
                            <input className="form-control" placeholder="Search for items" onChange={onSearch} />
                        </div>
                        <div className="ps-2 btn-group">
                            <Link className="btn btn-outline-primary  float-end" to="/items/create"><i className="fas fa-plus"></i></Link>
                            <Link className="btn btn-outline-primary" to="/items/report"><i className="fas fa-download"></i></Link>
                        </div>
                    </div>
                </div>
            </div>
            <Switch>
                <Route exact path={`${path}`}>
                    <Index />
                </Route>
                <Route path={`${path}/create`}>
                    <Add />
                </Route>
                <Route path={`${path}/qr/:id`}>
                    <QR />
                </Route>
                <Route path={`${path}/report`}>
                    <Report />
                </Route>
                <Route path={`${path}/withdraw/:id`}>
                    <Withdraw />
                </Route>
            </Switch>
        </div>
    );
}

interface ItemField {
    value: string,
    error: string | null
}

interface ItemForm {
    name: ItemField,
    desired: ItemField,
    actual: ItemField,
}

const required = (value: string): string | null => {
    var err = null;
    if (value === "") {
        err = "Required";
    }
    return err;
}

const isNumber = (value: string): string | null => {
    var pattern = /^\d+$/;
    var err = null;
    if (!pattern.test(value)) {
        err = "Requires an integer";
    }
    return err;
}

function Add() {
    const dispatch = useDispatch();
    const history = useHistory();
    const [allowCreate, setAllowCreate] = useState<boolean>(false);
    const [values, setValues] = useState<ItemForm>({
        name: { value: "", error: null },
        desired: { value: "", error: null },
        actual: { value: "", error: null },
    });

    useEffect(() => {
        if (values.name.value !== "" && values.desired.value !== "" && values.actual.value !== "" &&
            !values.name.error && !values.desired.error && !values.actual.error) {
            setAllowCreate(true);
        } else {
            setAllowCreate(false);
        }
    }, [values]);

    const handleChange = (fieldName: keyof ItemForm, validator?: (value: string) => string | null) => (
        e: React.ChangeEvent<HTMLInputElement>
    ) => {
        var err = null;
        if (validator) {
            err = validator(e.currentTarget.value)
        }        
        
        setValues({ ...values, [fieldName]: { value: e.currentTarget.value, error: err } });
    };

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        var item: Item = {
            name: values.name.value,
            desired: Number(values.desired.value),
            actual: Number(values.actual.value)
        }

        dispatch(createItem(item));
        history.push("/items");
    };

    return (
        <div className="row mt-4">
            <div className="col-12">
                <div className="card ">
                    <div className="card-body">
                        <h5 className="card-title">Create item</h5>
                        <form className="needs-validation">
                            <div className="row g-3">
                                <div className="col-12 col-sm-8">
                                    <label className="form-label">Name</label>
                                    <input className={`form-control ${values.name.error ? "is-invalid" : ""}`} value={values.name.value} onChange={handleChange("name", required)} />
                                    <div className="invalid-feedback">
                                        {values.name.error}
                                    </div>
                                </div>
                                <div className="col-6 col-sm-2">
                                    <label className="form-label">Desired</label>
                                    <input className={`form-control ${values.desired.error ? "is-invalid" : ""}`} value={values.desired.value} onChange={handleChange("desired", isNumber)} />
                                    <div className="invalid-feedback">
                                        {values.desired.error}
                                    </div>
                                </div>

                                <div className="col-6 col-sm-2">
                                    <label className="form-label">Actual</label>
                                    <input className={`form-control ${values.actual.error ? "is-invalid" : ""}`} value={values.actual.value} onChange={handleChange("actual", isNumber)} />
                                    <div className="invalid-feedback">
                                        {values.actual.error}
                                    </div>
                                </div>
                                <div className="col-12">
                                    <button type="submit" className={`btn btn-outline-primary float-end ${allowCreate ? "" : "disabled"}`} onClick={handleSubmit}>Create</button>
                                </div>
                            </div>
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

function Row(props: ItemProps) {
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

function Index() {
    const items = useSelector(selectItems);
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch(fetchItems(null));
    }, [dispatch]);

    return (
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
                            <Row key={item.id} item={item} />
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}

interface QRParams {
    id: string
}

function QR() {
    const { id } = useParams<QRParams>();
    const [item, setItem] = useState<Item>({
        name: "",
        actual: 0,
        desired: 0
    });

    useEffect(() => {
        axiosInstance.get("api/items/" + id).then((result: AxiosResponse<Item>) => {
            setItem(result.data);
        }).catch((error: AxiosError) => {
            console.log(error);
        });
    }, [id]);

    var QRCode = require('qrcode.react');

    return (
        <div className="row">
            <div className="col-12 col-sm-3">
                <div className="card mt-4">
                    <div className="card-body">
                        <h4 className="card-title">{item.name}</h4>
                        <div className="d-flex justify-content-center mt-4">
                            <QRCode className="d-flex justify-content-center" size={256} value={`${(window as any).STOQR_API_URL}items/withdraw/${item.id}`} />
                        </div>
                        <p className="float-end mt-2"><i>generated by STOQR</i></p>
                    </div>
                </div>
            </div>
        </div>
    );
}

function Report() {
    const items = useSelector(selectItems);

    return (
        <div className="row mt-4">
            <div className="col-12">
                <table className="table">
                    <thead>
                        <tr>
                            <th>#</th>
                            <th>Name</th>
                            <th>Needed</th>
                        </tr>
                    </thead>
                    <tbody className="align-middle">
                        {items.items.map((item: Item) => (
                            <tr key={item.id}>
                                <th>{item.id}</th>
                                <td>{item.name}</td>
                                <td>{item.desired - item.actual}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}

function Withdraw() {
    const dispatch = useDispatch();
    const { id } = useParams<QRParams>();
    const [item, setItem] = useState<Item>({
        name: "",
        actual: 0,
        desired: 0
    });
    useEffect(() => {
        axiosInstance.get("api/items/withdraw/" + id).then((result: AxiosResponse<Item>) => {
            setItem(result.data);
            dispatch(withdraw(result.data));
        }).catch((error) => {
            console.log(error);
        });
    }, [id, dispatch]);


    return (
        <div className="card mt-4">
            <div className="card-body">
                <h5 className="card-title">Stock updated succesfully!</h5>
                <h6 className="card-subtitle mb-2 text-muted">{item.name}</h6>
                <p className="card-text">Your stock is {item.actual}</p>
                <Link to="/items" className="card-link">View your stock</Link>
            </div>
        </div>
    );
}