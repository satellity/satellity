import React, {Component} from 'react';
import API from '../api/index.js';

class Index extends Component {
  constructor(props) {
    super(props);
    this.api = new API();

    let id = this.props.match.params.id;
    this.state = {
      group_id: id,
      messages: []
    }
  }

  componentDidMount() {
  }

  render() {
    return (
      <div>
        messages
      </div>
    )
  }
}

export default Index;
