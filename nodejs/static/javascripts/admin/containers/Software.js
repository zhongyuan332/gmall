import React, { Component } from 'react';
import { connect } from 'react-redux';

class Platform extends Component {
    render() {
        let { software } = this.props;
        return (
            <div className="app-footer">
                <div className="footer-content">
                    <span>{software.name} v{software.version}</span>
                    <span className="footer-divider">|</span>
                    <a href={software.officialURL} target="_blank">{software.officialURL}</a>
                </div>
            </div>
        );
    }
}

const mapStateToProps = (state) => {
    return {
        software: state.software
    };
}

export default connect(mapStateToProps)(Platform);