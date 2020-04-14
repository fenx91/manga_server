import React, { Component } from "react";
import ReactDOM from "react-dom";
import {BrowserRouter as Router, Switch, Route} from 'react-router-dom'
import Home from './components/home.js'
import MangaPage from './components/manga-page.js'
import "./index.css"

class App extends React.Component {
  render() {
    return (
    <Router>
      <Switch>
        <Route path="/" exact component={Home}/>
        <Route path="/mangapage/:mangaid" component={MangaPage}/>
      </Switch>
    </Router>
    );
  }
}
//{this.props.children}
const wrapper = document.getElementById("container");
wrapper ? ReactDOM.render(<App/>, wrapper) : false;

/*

    <Route path="home" component={Home} />
    <Route path="mangapage" component={MangaPage}/>
*/