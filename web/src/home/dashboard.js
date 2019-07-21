import style from './dashboard.scss';
import React, {Component} from 'react';
import {Link, Redirect} from 'react-router-dom';
import API from '../api/index.js';
import TopicItem from '../topics/item.js';
import GroupItem from '../groups/item.js';

class Dashboard extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    let user = this.api.user.readMe();
    this.state = {
      user: user,
      groups: [],
      topics: []
    }
  }

  componentDidMount() {
    if (!!this.state.user.user_id) {
      this.api.user.topics(this.state.user.user_id).then((data) => {
        this.setState({topics: data});
      });
      this.api.me.groups(3).then((data) => {
        this.setState({groups: data});
      });
    }
  }

  render() {
    let state = this.state;
    if (!state.user.user_id) {
      return (
        <Redirect to={{ pathname: "/" }} />
      )
    }

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
        <Link to='/user/groups' className={style.view}>{i18n.t('general.all')}</Link>
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
