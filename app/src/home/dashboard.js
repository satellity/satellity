import style from './dashboard.module.scss';
import React, {Component} from 'react';
import {Redirect} from 'react-router-dom';
import API from '../api/index.js';
import TopicItem from '../topics/item.js';
import Profile from '../users/profile.js';

class Dashboard extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    let user = this.api.user.local();
    this.state = {
      user: user,
      topics: []
    }
  }

  componentDidMount() {
    if (!!this.state.user.user_id) {
      this.api.user.topics(this.state.user.user_id).then((resp) => {
        if (resp.error) {
          return
        }
        this.setState({topics: resp.data});
      });
    }
  }

  render() {
    const i18n = window.i18n;
    let state = this.state;
    if (!state.user.user_id) {
      return (
        <Redirect to={{ pathname: "/" }} />
      )
    }

    let topics = state.topics.map((topic) => {
      return (
        <TopicItem topic={topic} key={topic.topic_id}/>
      )
    });

    return (
      <div className='container'>
        <main className='column main'>
          <div className={style.dashboard}>
            {
              state.topics.length !== 0 &&
              <div>
                <div className={style.title}>
                  {i18n.t('community.dashboard')}
                </div>
                <div className={style.section}>
                  {topics}
                </div>
              </div>
            }
          </div>
        </main>
        <aside className='column aside'>
          <Profile />
        </aside>
      </div>
    )
  }
}

export default Dashboard;
