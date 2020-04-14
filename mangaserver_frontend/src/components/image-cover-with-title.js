import React from "react"
import {Link} from "react-router-dom"
import "./image-cover-with-title.css"

class ImageCoverWithTitle extends React.Component {  // ICWT short for ImageCoverWithTitle.
    render() {
        return (
            <div className="ICWT">
                <Link to={this.props.linksto}>
                    <img className="imagecover" src={this.props.imgsrc} alt="image cover"></img>
                </Link>
                <div className="titlediv">
                    <p className="titletext">{this.props.title}</p>
                </div>
            </div>
        )
    }
}

export default ImageCoverWithTitle;