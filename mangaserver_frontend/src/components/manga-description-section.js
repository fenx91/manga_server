import React from "react"
import "./manga-description-section.css"

class MangaDescriptionSection extends React.Component {
    render() {
        return (
            <div className="manga-desc-sec">
                <img
                    className="feature-img"
                    style={{
                        width: this.props.imgwidth + 'px',
                        height: this.props.imgwidth * 2 / 3 + 'px',
                    }}
                    src={this.props.imgsrc}
                    alt="feature image"
                />                
                <div
                    className="desc-sec"
                    style={{
                        marginLeft: this.props.spacebetween + 'px',
                        marginRight: this.props.spacebetween + 'px',
                        marginTop: this.props.spacebetween * 0.5 + 'px',
                    }}
                >
                    <p
                        className="title-text"
                        style={{
                            fontSize: this.props.imgwidth / 15 + 'px',
                        }}
                    > {this.props.title}</p>
                    <p className="desc-text"> {this.props.desc} </p>
                </div>
            </div>
        )
    }
}

export default MangaDescriptionSection;