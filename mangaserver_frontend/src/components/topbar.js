import React from "react";
import "./topbar.css"
import {Link} from 'react-router-dom';

const Topbar = () => {
    return (
    <div className="header">
        <Link to="/">
            <img className="logo-with-title" src="/images/dragonblade-title-2.png" arl="logo"></img>
        </Link>
        <div>
            <ul className="nav-links">
                    <li><a href="#">Login</a></li>
                    <li><a href="#">Sign up</a></li>
            </ul>
        </div>
    </div>
    );
}
  
export default Topbar;
  