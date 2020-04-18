import React, { Component } from "react";
import Topbar from "./topbar.js"
import IcwtList from "./icwt-list.js";
import "./home.css"

class Home extends Component {
  constructor() {
    super();
    this.state = {
      icwt: []
    }
  }
  
  componentDidMount() {
    fetch(`/api/mangalist`)
      .then(res => res.json())
      .then(json => {
        const data = json.MangaDataList.map(mangadata => {
          return {
            key: mangadata.MangaId,
            imgsrc: "/static/manga/" + mangadata.MangaTitle + "/feature_img.jpg",
            title: mangadata.MangaTitle,
            link: "/mangapage/" + mangadata.MangaId,
          }
        });
        this.setState({
          icwt: data
        })
      });
  }

  render() {
    return (
      <div>
        <div> <Topbar/> </div>
        <div className="manga-cover-list">
           <div className="section-text"><p>All Mangas</p></div>
           <IcwtList icwtlist={this.state.icwt} imgheight="285px" imgwidth="435px" imgspace="35px" newTab={false}/>
        </div>
      </div>
    );
  }
}

export default Home;