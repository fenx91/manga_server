import React, { Component } from "react";
import ReactDOM from "react-dom";
import {BrowserRouter as Router, Switch, Route} from 'react-router-dom'
import Home from './components/home.js'
import MangaPage from './components/manga-page.js'
import Reader from './components/reader.js'
import NotFoundPage from './components/not-found-page.js'
import "./index.css"

class App extends React.Component {
  render() {
    return (
    <Router>
      <Switch>
        <Route path="/" exact component={Home}/>
        <Route path="/mangapage/:mangaid" component={MangaPage}/>
        <Route path="/reader/:mangaid/:chapterno" component={Reader}/>
        <Route component={NotFoundPage}/>
      </Switch>
    </Router>
    );
  }
}

const wrapper = document.getElementById("container");
wrapper ? ReactDOM.render(<App/>, wrapper) : false;