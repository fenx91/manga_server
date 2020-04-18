import React, { Component } from "react";
import Topbar from "./topbar.js"
import IcwtList from "./icwt-list.js";
import MangaDescriptionSection from "./manga-description-section.js"
import "./manga-page.css"

class MangaPage extends Component {
  constructor() {
    super();
    this.state = {
      icwt: [],
      mangatitle: "",
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
    fetch(`/api/mangainfo?mangaid=${this.props.match.params.mangaid}`)
    .then(res => res.json())
    .then((json) => {
      console.log(json);
      const icwt = [];
      let i = 1;
      for (i = 1; i <= json.ChapterCount; i++) {
        const formattedIndex = ("0" + i).slice(-2);
        icwt.push({
          key: i,
          imgsrc: "/static/manga/" + json.MangaTitle + "/" + formattedIndex + "/" + formattedIndex + "_000.jpg",
          title: "第" + formattedIndex + "卷",
          link: `/reader/${this.props.match.params.mangaid}/${i}`,
        })
      }
      this.setState({
        icwt: icwt,
        mangatitle: json.MangaTitle,
      });      
    });
    this.updateDimensions();
    window.addEventListener('resize', this.updateDimensions.bind(this));
  }

  render() {
    return (
      <div className="manga-page-div">
        <div> <Topbar/> </div>
        <div className="backgroundimage"></div>
        <div
          className="content-panel"
          style={{
            width: this.state.w * 0.6,
            left: this.state.w * 0.2,
            top: this.state.h * 0.05,
            position: 'relative',
          }}
        > 
          <div
            className="manga-desc-sec-div"
            style={{
              marginLeft: (this.state.w - 4) * 0.6 * 0.25 / 10 * 2,
              marginTop: (this.state.w - 4) * 0.6 * 0.25 / 10 * 2,
            }}
          >
            <MangaDescriptionSection
              imgsrc={this.state.mangatitle == "" ? "" : "/static/manga/" + this.state.mangatitle + "/feature_img.jpg"}
              title={this.state.mangatitle}
              desc=""
              imgwidth={(this.state.w - 4) * 0.6 * 0.75 / 4 * 2 + (this.state.w - 4) * 0.6 * 0.25 / 10 * 2}
              spacebetween={(this.state.w - 4) * 0.6 * 0.25 / 10 * 2}
            />
          </div>
          <div
            className="chapter-cover-list"
            style={{
              marginLeft: (this.state.w - 4) * 0.6 * 0.25 / 10,
            }}
          >
            <IcwtList
              icwtlist={this.state.icwt}
              imgwidth={(this.state.w - 4) * 0.6 * 0.75 / 4}
              imgheight={(this.state.w - 4) * 0.6 * 0.75 / 4 * 3 / 2}
              imgspace={(this.state.w - 4) * 0.6 * 0.25 / 10}
              marginTop={(this.state.w - 4) * 0.6 * 0.25 / 10 * 2}
              textSize={Math.round((this.state.w - 4) * 0.6 * 0.75 / 4 * 20 / 200)}
              newTab={true}
            />
          </div>          
        </div>
      </div>
    );
  }
}

export default MangaPage;