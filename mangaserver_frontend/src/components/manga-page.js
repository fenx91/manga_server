import React, { Component } from "react";
import Topbar from "./topbar.js"
import IcwtList from "./icwt-list.js";
import "./manga-page.css"

class MangaPage extends Component {
  constructor() {
    super();
    this.state = {
      icwt: []
    }
  }
  
  componentDidMount() {
    fetch(`http://localhost:80/api/mangainfo?mangaid=${this.props.match.params.mangaid}`)
    .then(res => res.json())
    .then(json => {
      console.log(json);
      const icwt = [];
      let i = 1;
      for (i = 1; i < json.ChapterCount; i++) {
        const formattedIndex = ("0" + i).slice(-2);
        icwt.push({
          key: i,
          imgsrc: "/static/manga/" + json.MangaTitle + "/" + formattedIndex + "/" + formattedIndex + "_000.jpg",
          title: "第" + formattedIndex + "卷",
          link: "/",
        })
      }
      this.setState({
        icwt: icwt,
      });
      console.log(icwt);
    });
  }

  render() {
    return (
      <div>
        <div> <Topbar/> </div>
        <div className="chapter-cover-list"> <IcwtList icwtlist={this.state.icwt}/> </div>
      </div>
    );
  }
}

export default MangaPage;