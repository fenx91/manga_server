import React from 'react'
import './tool-button.css'

class ToolButton extends React.Component{
    constructor() {
        super();
    }

    render() {
        return (
            <div className="tools-button-background"
                style={{
                    height: this.props.h + 'px',
                    width: this.props.h + 'px',
                    position: 'absolute',
                    left: this.props.left + 'px',
                    top: this.props.top + 'px',
                }}
                onClick={this.props.clickHandler}
            >
                <img
                className="tools-button"
                src="/images/tools-button.png"
                style={{
                    height: this.props.h + 'px',
                    width: this.props.h + 'px',
                    position: 'absolute',
                }}
                />
            </div>
        );  
    }
}

export default ToolButton;