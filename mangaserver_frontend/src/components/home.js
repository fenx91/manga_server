import React, { Component } from "react";
import Topbar from "./topbar.js"
import IcwtList from "./icwt-list.js";
import "./home.css"

class Home extends Component {
  constructor() {
    super();
    this.state = {
      icwt: [],
      w: 400,
      h: 400,
    }
  }

  updateDimensions() {
    this.setState({
      w: window.innerWidth,
      h: window.innerHeight,
    })
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
    this.updateDimensions();
    window.addEventListener('resize', this.updateDimensions.bind(this));
  }

  render() {
    return (
      <div> 
        <div> <Topbar/> </div>
        <div
          className="manga-cover-list"
          style={{
            width: this.state.w * 0.6,
            left: this.state.w * 0.2,
            top: this.state.h * 0.05,
            position: 'relative',
          }}
        >
          <div
            className="section-text"
            style={{
              marginLeft: ((this.state.w - 4)* 0.6 * 0.2) / 4 + 'px',
              marginTop: ((this.state.w - 4)* 0.6 * 0.2) / 4 + 'px',
              fontSize: Math.round(this.state.w / 40) + 'px'
            }}
          >
            <p>All Mangas</p>
          </div>
          <IcwtList
            icwtlist={this.state.icwt}
            imgwidth={(this.state.w - 4) * 0.6 / 2 * 0.8}
            imgheight={((this.state.w - 4) * 0.6 / 2 * 0.8) * 2 / 3 }
            imgspace={((this.state.w - 4)* 0.6 * 0.2) / 4}
            marginTop={((this.state.w - 4)* 0.6 * 0.2) / 4}
            textSize={Math.round((this.state.w - 4) * 0.6 / 2 * 0.8 * 20 / 460)}
            newTab={false}
          />
        </div>
      </div>
    );
  }
}

export default Home;