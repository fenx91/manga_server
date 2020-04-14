import React from "react"
import ImageCoverWithTitle from "./image-cover-with-title.js"
import "./icwt-list.css"

class IcwtList extends React.Component {  // ICWT short for ImageCoverWithTitle.
    render() {
        console.log("-" + this.props.imgspace)
        return (
            <div className="icwt-list-div" style={{"margin-left": "-" + this.props.imgspace, "margin-right": ("-" + this.props.imgspace)}}>
                    {this.props.icwtlist.map((icwt, index) => {
                        return (
                            <div className="list-item" key={icwt.key} style={{"margin-left": this.props.imgspace, "margin-right": this.props.imgspace}}>
                                <ImageCoverWithTitle
                                    imgsrc={icwt.imgsrc}
                                    title={icwt.title}
                                    linksto={icwt.link}
                                    imgwidth={this.props.imgwidth}
                                    imgheight={this.props.imgheight}
                                />
                            </div>
                        )
                    })}
            </div>
        )
    }
}

export default IcwtList;