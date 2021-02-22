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
                <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarToggler"
                    aria-controls="navbarToggler" aria-expanded="false" aria-label="Toggle navigation">
                    <span className="navbar-toggler-icon"></span>
                </button>
                <div className="collapse navbar-collapse" id="navbarToggler">
                    <ul className="navbar-nav me-auto mb-2 mb-lg-0"></ul>                    
                    <Link className="btn btn-link" to="items/create"><i className="fas fa-plus"></i></Link>
                    <Link className="btn btn-link" to="/items/report"><i className="fas fa-download"></i></Link>
                    {itemsWithNoStock > 0 && <Link className="btn btn-link" to="/items/report"><span className="badge bg-danger">{itemsWithNoStock}</span></Link>}                        
                </div>
            </div>
        </nav>
    );
}