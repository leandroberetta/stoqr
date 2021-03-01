import React from 'react';
import { Link } from "react-router-dom";
import { selectItems } from '../store/itemsSlice';
import { useSelector } from 'react-redux';
import { Item } from '../model/item';

export function Menu() {
    const items = useSelector(selectItems);
    var itemsWithNoStock = 0;
     
    items.items.forEach((item: Item) => {
        if (item.actual === 0) {
            itemsWithNoStock = itemsWithNoStock + 1;
        }
    });

    return (
        <nav className="navbar navbar-expand-lg navbar-light stoqr-navbar">
            <div className="container">
                <Link className="navbar-brand" to="/">STOQR</Link>
                {itemsWithNoStock > 0 && <Link className="btn btn-link" to="/items/report"><span className="badge bg-danger">{itemsWithNoStock}</span></Link>}                                        
            </div>
        </nav>
    );
}