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
    }
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
  }

  render() {
    return (
      <div className="manga-page-div">
        <div> <Topbar/> </div>
        <div className="backgroundimage"></div>
        <div className="content-panel"> 
          <div className="manga-desc-sec-div">
            <MangaDescriptionSection
              //imgsrc="/images/feature_img.jpg"
              imgsrc={this.state.mangatitle == "" ? "" : "/static/manga/" + this.state.mangatitle + "/feature_img.jpg"}
              //title="ぼくたちは勉強ができない"
              title={this.state.mangatitle}
              desc="Some description..."
            />
          </div>
          <div className="chapter-cover-list">
            <IcwtList icwtlist={this.state.icwt} imgwidth="200px" imgheight="300px" imgspace="20px" newTab={true}/>
          </div>          
        </div>
      </div>
    );
  }
}

export default MangaPage;