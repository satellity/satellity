import style from './index.module.scss';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../../api/index.js';

class AdminTopic extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {topics: []};
    this.handleDelete = this.handleDelete.bind(this);
  }

  componentDidMount() {
    this.api.topic.admin.index().then((resp) => {
      if (resp.error) {
        return
      }
      this.setState({topics: resp.data});
    });
  }

  handleDelete(e, id, title) {
    e.preventDefault();
    let c = window.confirm(`Delete: ${title}`);
    if (c) {
      this.api.topic.admin.delete(id).then((resp) => {
        if (resp.error) {
          return
        }

        let topics = this.state.topics.filter((topic) => {
          return topic.topic_id !== id;
        });
        this.setState({topics: topics});
      });
    }
  }

  render() {
    const state = this.state;

    const listTopics = state.topics.map((topic) => {
      return (
        <li key={topic.topic_id}>
            {topic.user.nickname} |
          <Link to={`/topics/${topic.topic_id}`}>{topic.title}</Link>
          <div className={style.time}>
              {topic.topic_id} | {topic.created_at} |
            <Link to={`/topics/${topic.topic_id}/edit`} >EIDT</Link> |
            <Link to='' onClick={(e) => this.handleDelete(e, topic.topic_id, topic.title)} >DELETE</Link>
          </div>
        </li>
      )
    });

    return (
      <div>
        <h1 className='welcome'>
            Here is the list of  topics.
        </h1>
        <div className='panel'>
          <ul>
              {listTopics}
          </ul>
        </div>
      </div>
    );
  }
}

export default AdminTopic;
