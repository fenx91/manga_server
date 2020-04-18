import React from "react";
import ToolBar from './tool-bar.js'
import ToolButton from './tool-button.js'
import ChaptersList from './chapters-list.js'
import "./reader.css"

class Reader extends React.Component {
    constructor() {
        super();
        this.state = {
            w: 800,
            h: 600,
            l_w: 400,
            l_h: 600,
            l_top: 0,
            l_left: 0,
            r_w: 400,
            r_h: 600,
            r_top: 0,
            r_left: 0,
            l_imgsrc: "",
            r_imgsrc: "",
            lpic_wh_ratio:0.67,
            rpic_wh_ratio:0.67,
            rightPageAt: 0,
            chapterPageCount: 0,
            totalChapter: 0,
            enableChaptersList: false,
            enableToolBar: false,
        }
        // reader stats.
        this.mangaTitle = "";
        this.mangaId = 0;
        this.currentChapter = 0;
        this.chapterPageCount = 0;        
        this.rightPageAt = 0;
        this.totalChapter = 0;
        this.mangaInfoJson = {};
        this.dataReady = false;
        // timer that handles the timeout to disbale toolbar
        this.disableToolBarTimer = null;
        // indicates whether img loading has finished.
        this.lImgLoaded = false;
        this.rImgLoaded = false;
        // constants
        this.timeToDisableToolBar = 2000;  // ms
    }

    init() {
        this.mangaTitle = this.mangaInfoJson.MangaTitle;
        this.totalChapter = this.mangaInfoJson.MangaChapterInfoList.length;
        let dataReady = true, chapterNoSum = 0;
        for (let i = 0; i < this.totalChapter; ++i) {
            chapterNoSum += this.mangaInfoJson.MangaChapterInfoList[i].ChapterNo;
            if (this.mangaInfoJson.MangaChapterInfoList[i].PageCount <= 0) {
                dataReady = false;
            }
        }
        if (chapterNoSum != (this.totalChapter * (this.totalChapter + 1)) / 2) dataReady = false;
        this.dataReady = dataReady;
        if (!this.dataReady) {
            console.warn("something wrong with the retrieved chapter info list.");
            return
        } 
        this.chapterPageCount = this.getChapterPageCount(this.currentChapter);
    }

    getChapterPageCount(chapterNo) {
        for (let i = 0; i < this.totalChapter; ++i) {
            if (this.mangaInfoJson.MangaChapterInfoList[i].ChapterNo == chapterNo) {
                return this.mangaInfoJson.MangaChapterInfoList[i].PageCount;
            }
        }
        return 0;
    }

    // this.setState() called in this func.
    updateDimensions() {
        const w = window.innerWidth;
        const h = window.innerHeight;
        const halfScreenWhRatio = w / 2 / h;
        let new_l_w, new_l_h, new_r_w, new_r_h = 0;
        let new_l_top, new_l_left, new_r_top, new_r_left = 0;
        if (halfScreenWhRatio > this.state.lpic_wh_ratio) {
            new_l_h = h;
            new_l_w = this.state.lpic_wh_ratio * new_l_h;
            new_l_top = 0;
            new_l_left = Math.floor(w / 2) - new_l_w;
        } else {
            new_l_w = w * (this.state.lpic_wh_ratio) / (this.state.lpic_wh_ratio + this.state.rpic_wh_ratio);
            new_l_h = new_l_w / this.state.lpic_wh_ratio;
            new_l_left = 0;
            new_l_top = (h - new_l_h) / 2;
        }

        if (halfScreenWhRatio > this.state.rpic_wh_ratio) {
            new_r_h = h;
            new_r_w = this.state.rpic_wh_ratio * new_r_h;
            new_r_top = 0;
            new_r_left = new_l_left + new_l_w;
        } else {
            new_r_w = w - new_l_w;
            new_r_h = new_r_w / this.state.rpic_wh_ratio;
            new_r_left = w * this.state.lpic_wh_ratio / (this.state.lpic_wh_ratio + this.state.rpic_wh_ratio);
            new_r_top = (h - new_r_h) / 2;
        }

        this.setState({
            w: w,
            h: h,
            l_w: new_l_w,
            l_h: new_l_h,
            l_top: new_l_top,
            l_left: new_l_left,
            r_w: new_r_w,
            r_h: new_r_h,
            r_top: new_r_top,
            r_left: new_r_left,
        });
    }

    getLeftImageSrc() {
        const leftPageAt = this.rightPageAt + 1;
        if (leftPageAt < this.chapterPageCount) {
            const formattedCurrentChapter = ("0" + this.currentChapter).slice(-2);
            const formattedPageNo = ("00" + leftPageAt).slice(-3);
            return `/static/manga/${this.mangaTitle}/${formattedCurrentChapter}/${formattedCurrentChapter}_${formattedPageNo}.jpg`;
        } else if (this.currentChapter < this.totalChapter){
            const formattedNextChapter = ("0" + (this.currentChapter + 1)).slice(-2);
            const formattedPageNo = ("00" + (leftPageAt - this.chapterPageCount)).slice(-3);
            return `/static/manga/${this.mangaTitle}/${formattedNextChapter}/${formattedNextChapter}_${formattedPageNo}.jpg`;
        } else {
            return ""
        }
    }

    getRightImageSrc() {
        if (this.rightPageAt >= 0 && this.rightPageAt < this.chapterPageCount) {
            const formattedCurrentChapter = ("0" + this.currentChapter).slice(-2);
            const formattedPageNo = ("00" + this.rightPageAt).slice(-3);
            return `/static/manga/${this.mangaTitle}/${formattedCurrentChapter}/${formattedCurrentChapter}_${formattedPageNo}.jpg`;
        } else {
            return "";
        }
    }
 
    // this.setState() called in this func.
    setPicSrc() {
        this.setState({
            rightPageAt: this.rightPageAt,
            chapterPageCount: this.chapterPageCount,
        });
        // Load the left side pic if needed.
        const rsrc = this.getRightImageSrc();
        console.log("rpic src: " + rsrc);
        if (rsrc) {
            this.rImgLoaded = false;
            const rPicImg = new Image();
            rPicImg.onload = (function(reader){
                return function() {
                    reader.setState({
                        r_imgsrc: this.src,
                        rpic_wh_ratio: this.width / this.height,
                    });
                    reader.rImgLoaded = true;
                    if (reader.rImgLoaded && reader.lImgLoaded) {
                        reader.updateDimensions();
                    }
                };
            })(this);
            rPicImg.src = rsrc;
        } else {
            this.rImgLoaded = true;
            this.setState({
                r_imgsrc: "",
            })
        }
        // Load the left side pic if needed.
        const lsrc = this.getLeftImageSrc();
        console.log("lpic src: " + lsrc);
        if (lsrc) {
            this.lImgLoaded = false;
            const lPicImg = new Image();        
            lPicImg.onload = (function(reader){
                return function() {
                    reader.setState({
                        l_imgsrc: this.src,
                        lpic_wh_ratio: this.width / this.height,
                    });
                    reader.lImgLoaded = true;
                    if (reader.rImgLoaded && reader.lImgLoaded) {
                        reader.updateDimensions();
                    }
                };
            })(this);
            lPicImg.src = lsrc;   
        } else {
            this.lImgLoaded = true;
            this.setState({
                l_imgsrc: "",
            })
        }          
    }

    componentDidMount() {
        this.setState({
            w: window.innerWidth,
            h: window.innerHeight,
        });
        this.currentChapter = parseInt(this.props.match.params.chapterno);
        this.mangaId = parseInt(this.props.match.params.mangaid);
        fetch(`/api/chapterpagecount?mangaid=${this.mangaId}`)
        .then(res => res.json())
        .then(json => {
            this.mangaInfoJson = json;
            this.init();

            if (this.dataReady) {
                this.setState({
                    totalChapter: this.totalChapter,
                });
                this.handleToolButtonClick(null);
                this.setPicSrc();
            }
        });
        window.addEventListener("resize", this.updateDimensions.bind(this));
        this.lcrElem.addEventListener("wheel", this.handleWheel.bind(this), {passive: false});
        this.rcrElem.addEventListener("wheel", this.handleWheel.bind(this), {passive: false});
    }

    nextPage() {
        if (this.rightPageAt + 2 < this.chapterPageCount) {
            this.rightPageAt += 2;
            this.setPicSrc();
        } else if (this.currentChapter < this.totalChapter) {  // goes to next chapter if possible
            this.rightPageAt = this.rightPageAt + 2 - this.chapterPageCount;
            this.currentChapter++;
            this.chapterPageCount = this.getChapterPageCount(this.currentChapter);
            this.setPicSrc();
            window.history.pushState({},"", `/reader/${this.mangaId}/${this.currentChapter}`)
        } else {
            // otherwise do nothing.
        }
    }

    prevPage() {
        if (this.rightPageAt >= 2) {
            this.rightPageAt -= 2;
            this.setPicSrc();
        } else if (this.currentChapter > 1) {  // go to prev chapter if possible
            this.currentChapter--;
            this.chapterPageCount = this.getChapterPageCount(this.currentChapter);
            this.rightPageAt = this.rightPageAt + this.chapterPageCount - 2;
            this.setPicSrc();
            window.history.pushState({},"", `/reader/${this.mangaId}/${this.currentChapter}`)
        } else if (this.currentChapter == 1 && this.rightPageAt == 1) { // only display left pic.
            this.rightPageAt -= 2;
            this.setPicSrc();
        } else {
            // otherwise do nothing.
        }
    }

    handleLeftClick(e) {
        e.preventDefault();
        this.disableToolBar();
        this.nextPage();
    }

    handleRightClick(e) {
        e.preventDefault();
        this.disableToolBar();
        this.prevPage();    
    }

    handleWheel(event) {
        event.preventDefault();
        if (event.deltaY > 0) {
            this.nextPage();
        } else {
            this.prevPage();
        }
    }

    handleDragPageBar(pageNo) {
        this.rightPageAt = pageNo - 1;
        this.setPicSrc();
    }

    handleToolButtonClick(e) {   
        if (this.disableToolBarTimer) {
            window.clearTimeout(this.disableToolBarTimer);
            this.disableToolBarTimer = null;
        }
        const enableToolBar = !this.state.enableToolBar;
        if (enableToolBar) {
            this.setState({
                enableToolBar: enableToolBar,
            });
            this.disableToolBarTimer = window.setTimeout(function(reader) {
                return function() {
                    reader.disableToolBar();
                }
            }(this), this.timeToDisableToolBar);
        } else {
            this.disableToolBar();
        }
    }

    handleChapterListClick(e) {
        e.preventDefault();
        const enableChaptersList = !this.state.enableChaptersList;
        this.setState({
            enableChaptersList: enableChaptersList,
        })
    }

    handleMouseEnterToolBar(e) {
        e.preventDefault();
        if (this.disableToolBarTimer) {
            window.clearTimeout(this.disableToolBarTimer);
            this.disableToolBarTimer = null;
        }
    }

    handleMouseLeaveToolBar(e) {
        e.preventDefault();
        if (this.disableToolBarTimer) {
            window.clearTimeout(this.disableToolBarTimer);
            this.disableToolBarTimer = null;
        }
        this.disableToolBarTimer = window.setTimeout(function(reader) {
            return function() {
                reader.disableToolBar();
            }
        }(this), this.timeToDisableToolBar);
    }

    disableToolBar() {
        this.setState({
            enableToolBar: false,
            enableChaptersList: false,
        });
    }

    render() {
        return (
            <div className="reader-background" style={{width:this.state.w, height:this.state.h}} >
                <div
                    ref={elem => this.lcrElem = elem}
                    className="left-click-receiver"
                    style={{
                        width: this.state.w / 2 + "px",
                        height: this.state.h + "px",
                        position: "absolute",
                        top: "0px",
                        left: "0px",
                    }}
                    onClick={(e) => this.handleLeftClick(e)}
                />
                <div
                    ref={elem => this.rcrElem = elem}
                    className="right-click-receiver"
                    style={{
                        width: this.state.w / 2 + "px",
                        height: this.state.h + "px",
                        position: "absolute",
                        top: "0px",
                        left: this.state.w / 2 + "px",
                    }}
                    onClick={(e) => this.handleRightClick(e)}
                />
                <img
                    src={this.state.l_imgsrc}
                    className="left-canvas"
                    style={{
                        width: this.state.l_w + "px",
                        height: this.state.l_h + "px",
                        position: "absolute",
                        top:this.state.l_top + "px",
                        left:this.state.l_left + "px",
                        opacity: this.state.l_imgsrc ? "1" : "0",
                    }}
                />
                <img
                    src={this.state.r_imgsrc}
                    className="right-canvas"
                    style={{
                        width: this.state.r_w + "px",
                        height: this.state.r_h + "px",
                        position: "absolute",
                        top: this.state.r_top + "px",
                        left: this.state.r_left + "px",
                        opacity: this.state.r_imgsrc ? "1" : "0",
                    }}
                />

                <ToolButton
                    h={this.state.h / 15}
                    left={this.state.w - this.state.h / 15 - this.state.w / 15}
                    top={this.state.h * 14 / 15 - this.state.h / 15 - this.state.h / 20}
                    clickHandler={(e) => {this.handleToolButtonClick(e)}}
                />

                <ToolBar
                    w={this.state.w}
                    h={this.state.h / 15}
                    left={0}
                    top={this.state.enableToolBar ? this.state.h * 14 / 15 : this.state.h}
                    horizontalPadding={(this.state.w - this.state.l_w - this.state.r_w) / 2}
                    chapterPageCount={this.state.chapterPageCount}
                    currentPageNo={parseInt(this.state.rightPageAt) + 1} 
                    prevChapterLink={this.currentChapter > 1 ? `/reader/${this.mangaId}/${this.currentChapter - 1}` : null}
                    nextChapterLink={this.currentChapter < this.totalChapter ? `/reader/${this.mangaId}/${this.currentChapter + 1}` : null}
                    changeHandler={(e) => {this.handleDragPageBar(e.target.value)}}
                    enabled={this.state.enableToolBar}
                    chapterListClickHandler={(e) => {this.handleChapterListClick(e)}}
                    mouseEnterHandler={(e) => {this.handleMouseEnterToolBar(e)}}
                    mouseLeaveHandler={(e) => {this.handleMouseLeaveToolBar(e)}}
                />

                <ChaptersList
                    w={(this.state.l_w + this.state.r_w) / 8 * 1.25}
                    h={this.state.h / 4}
                    left={(this.state.w - this.state.l_w - this.state.r_w) / 2}
                    top={(this.state.h * 14 / 15) - (this.state.h / 4 * 9 / 8)}
                    mangaId={this.mangaId}
                    totalChapter={this.state.totalChapter}
                    enabled={this.state.enableChaptersList}
                    mouseEnterHandler={(e) => {this.handleMouseEnterToolBar(e)}}
                    mouseLeaveHandler={(e) => {this.handleMouseLeaveToolBar(e)}}
                />
            </div>
        );
    }
}

export default Reader;