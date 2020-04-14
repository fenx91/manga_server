import React from "react"
import ImageCoverWithTitle from "./image-cover-with-title.js"
import "./icwt-list.css"

class IcwtList extends React.Component {  // ICWT short for ImageCoverWithTitle.
    render() {
        return (
            <div className="icwt-list-div">
                    {this.props.icwtlist.map((icwt, index) => {
                        return (
                            <div className="list-item" key={icwt.key}>
                                <ImageCoverWithTitle
                                    imgsrc={icwt.imgsrc}
                                    title={icwt.title}
                                    linksto={icwt.link}
                                />
                            </div>
                        )
                    })}
            </div>
        )
    }
}

export default IcwtList;