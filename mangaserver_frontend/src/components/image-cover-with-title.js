import React from "react"
import {Link} from "react-router-dom"
import "./image-cover-with-title.css"

class ImageCoverWithTitle extends React.Component {  // ICWT short for ImageCoverWithTitle.
    render() {
        return (
            <div className="ICWT">
                <Link to={this.props.linksto} target={this.props.newTab ? "_blank" : ""}>
                    <img
                        className="imagecover"
                        src={this.props.imgsrc}
                        alt="image cover"
                        style={{
                            width: this.props.imgwidth + 'px',
                            height: this.props.imgheight + 'px',
                        }}
                    />
                </Link>
                <div
                    className="titlediv"
                    style={{
                        width: this.props.imgwidth + 'px',
                        fontSize: this.props.textSize + 'px',
                    }}
                >
                    <p
                        className="titletext"
                    >
                        {this.props.title}
                    </p>
                </div>
            </div>
        )
    }
}

export default ImageCoverWithTitle;