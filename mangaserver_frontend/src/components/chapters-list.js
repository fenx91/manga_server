import React from 'react' 
import './chapters-list.css'

class ChaptersList extends React.Component {
    constructor() {
        super()
        this.state = {
            chapterNo: [],
        };
    }

    componentWillReceiveProps(newProps) {
        const chapterNo = new Array();
        for (let i = 1; i <= newProps.totalChapter; ++i) {
           chapterNo.push(i);
        }
        this.setState({
            chapterNo: chapterNo,
        });
    }

    render() {
        return (
            <div className="overall-chapters-list"
                style={{
                    width: this.props.w,
                    position: 'absolute',
                    top: this.props.top,
                    left: this.props.left,
                    display: this.props.enabled ? '' : 'none',
                }}
                onMouseEnter={this.props.mouseEnterHandler}
                onMouseLeave={this.props.mouseLeaveHandler}
            >
                <h1 className="chapter-list-title"
                    style={{
                        height: this.props.h / 8,
                        width: this.props.w,
                    }}
                >All Chapters</h1>
                <ul className="chapter-list-items"
                    style={{
                        height: this.props.h,
                    }}
                >
                { this.state.chapterNo.map((element) => {
                    const formatted = ("0" + element).slice(-2);
                    return (
                        <li key={element}>
                            <a className="chapter-item" href={`/reader/${this.props.mangaId}/${element}`}> 第 {formatted} 卷 </a>
                        </li>
                    ); 
                }) }
                </ul>
            
            </div>
        );
    }
}

export default ChaptersList;