import style from './dashboard.scss';
import React, {Component} from 'react';
import {Link} from 'react-router-dom';
import API from '../api/index.js';
import TopicItem from '../topics/item.js';
import GroupItem from '../groups/item.js';

class Dashboard extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {
      groups: [],
      topics: []
    }
  }

  componentDidMount() {
    let user = this.api.user.readMe();
    this.api.user.topics(user.user_id).then((data) => {
      this.setState({topics: data});
    });
    this.api.me.groups(3).then((data) => {
      this.setState({groups: data});
    });
  }

  render() {
    let state = this.state;
    const groups = state.groups.map((group) => {
      return (
        <div key={group.group_id} className={style.group}>
          <GroupItem group={group} />
        </div>
      )
    })
    const topics = state.topics.map((topic) => {
      return (
        <TopicItem topic={topic} key={topic.topic_id}/>
      )
    });

    return (
      <div className={style.dashboard}>
        {i18n.t('group.dashboard')}
        <div className={style.groups}>
          {groups}
        </div>

        {i18n.t('community.dashboard')}
        <div className={style.section}>
          {topics}
        </div>
      </div>
    )
  }
}

export default Dashboard;
