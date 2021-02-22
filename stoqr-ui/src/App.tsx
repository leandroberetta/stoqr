import React from 'react';
import { Switch, Route, Redirect } from "react-router-dom";
import './App.css';
import { Items } from './components/Items'
import { Menu } from './components/Menu'

function App() {
  return (
    <div>
      <Menu />
      <div className="container">
        <Switch>
          <Route exact path="/">
            <Redirect to="/items" />
          </Route>
          <Route path="/items">
            <Items />
          </Route>
        </Switch>
      </div>
    </div>
  );
}

export default App;
