import React, {Component} from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';

class CommentNew extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.state = {topic_id: props.topicId, body: ''};
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleSubmit(e) {
    e.preventDefault();
    this.api.comment.create(this.state, () => {
      this.setState({body: ''});
    });
  }

  render() {
    if (!this.api.user.loggedIn()) {
      return ''
    }
    return (
      <View onSubmit={this.handleSubmit} onChange={this.handleChange} state={this.state} />
    )
  }
}

const View = ({onSubmit, onChange, state}) => {
  return (
    <form onSubmit={(e) => onSubmit(e)}>
      <input type='hidden' name='topic_id' defaultValue={state.topic_id} />
      <div>
        <textarea type='text' name='body' value={state.body} onChange={(e) => onChange(e)} />
      </div>
      <div className='action'>
        <input type='submit' value='SUBMIT' />
      </div>
    </form>
  )
};

export default CommentNew;
