import style from './dashboard.scss';
import React, {Component} from 'react';
import {Link, Redirect} from 'react-router-dom';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import API from '../api/index.js';
import TopicItem from '../topics/item.js';
import GroupItem from '../groups/item.js';

class Dashboard extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    let user = this.api.user.local();
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

    let groups = state.groups.map((group) => {
      return (
        <div key={group.group_id} className={style.group}>
          <GroupItem group={group} to='messages' />
        </div>
      )
    })
    let topics = state.topics.map((topic) => {
      return (
        <TopicItem topic={topic} key={topic.topic_id}/>
      )
    });

    return (
      <div className={style.dashboard}>
        <div className={style.create}>
          {
            state.groups.length < 3 &&
            <Link to='/groups/new'>{i18n.t('group.new')}</Link>
          }
        </div>
        {
          state.groups.length != 0 &&
          <div>
            <div className={style.head}>
              <div className={style.name}>
                {i18n.t('group.dashboard')}
              </div>
              <Link to='/user/groups' className={style.view}>{i18n.t('general.all')}</Link>
            </div>
            <div className={style.groups}>
              {groups}
            </div>
          </div>
        }

        {
          state.topics.length != 0 &&
          <div>
            <div className={style.title}>
                {i18n.t('community.dashboard')}
              <Link to='/topics/new'>
                <FontAwesomeIcon icon={['fa', 'plus']} />
              </Link>
            </div>
            <div className={style.section}>
              {topics}
            </div>
          </div>
        }
      </div>
    )
  }
}

export default Dashboard;
