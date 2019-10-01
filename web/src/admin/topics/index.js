import style from './index.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../../api/index.js';

class AdminTopic extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {topics: []};
  }

  componentDidMount() {
    this.api.topic.admin.index().then((resp) => {
      if (resp.error) {
        return
      }
      this.setState({topics: resp.data});
    });
  }

  render() {
    return (
      <TopicIndex topics={this.state.topics} />
    )
  }
}

const TopicIndex = (props) => {
  const listTopics = props.topics.map((topic) => {
    return (
      <li key={topic.topic_id}>
        {topic.user.nickname} | {topic.title}
        <div className={style.time}>
          {topic.topic_id} | {topic.created_at}
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
  )
}

export default AdminTopic;
