import React from 'react'
import { Link } from 'react-router-dom'
import './tool-bar.css'

class ToolBar extends React.Component {
    constructor() {
        super();
    }

    render() {
        return (
            <div
                className="overall-tool-bar"
                style={{
                    height: this.props.h + 'px',
                    width: this.props.w  + 'px',
                    left: this.props.left + 'px',
                    top: this. props.top + 'px',
                    postion: 'absolute',
                    paddingLeft: this.props.horizontalPadding + 'px',
                    paddingRight: this.props.horizontalPadding + 'px',
                    
                }}
                onMouseEnter={this.props.mouseEnterHandler}
                onMouseLeave={this.props.mouseLeaveHandler}
            >
                <p className="chapter-list toolbar-text" onClick={this.props.chapterListClickHandler}>Chapter list</p>
                <p className="page-number toolbar-text"> {this.props.currentPageNo} / {this.props.chapterPageCount} </p>
                <input
                    type="range" id="page-input-range" name="page-range"
                    min="1"
                    max={this.props.chapterPageCount}
                    step="2"
                    onInput={this.props.changeHandler}
                />
                <a
                    className={`toolbar-text${this.props.prevChapterLink ? " toolbar-link" : ""}`}
                    href={this.props.prevChapterLink}
                    style={{
                        color: `${this.props.prevChapterLink ? "" : "grey"}`,
                    }}
                >
                    <p>Prev</p>
                </a>
                <a 
                    className={`toolbar-text${this.props.nextChapterLink ? " toolbar-link" : ""}`}
                    href={this.props.nextChapterLink}
                    style={{
                        color: `${this.props.nextChapterLink ? "" : "grey"}`,
                    }}
                >
                    <p>Next</p>
                </a>
            </div>
        );
    }
}

export default ToolBar;