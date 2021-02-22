import React from 'react';
import { useSelector } from 'react-redux';
import { selectItems } from '../store/itemsSlice';
import { Item } from '../model/item';

export function Report() {
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
                                <th>{ item.id }</th>
                                <td>{ item.name}</td>
                                <td>{ item.desired - item.actual }</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}


