import React from "react"
import "./manga-description-section.css"

class MangaDescriptionSection extends React.Component {
    render() {
        return (
            <div className="manga-desc-sec">
                <img className="feature-img" src={this.props.imgsrc} alt="feature image"/>                
                <div className="desc-sec">
                    <p className="title-text"> {this.props.title}</p>
                    <p className="desc-text"> {this.props.desc} </p>
                </div>
            </div>
        )
    }
}

export default MangaDescriptionSection;