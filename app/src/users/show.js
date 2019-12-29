import style from './show.module.scss';
import moment from 'moment';
import React, { Component } from 'react';
import API from '../api/index.js';
import TopicItem from '../topics/item.js';

class Show extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {
      id: props.match.params.id,
      user: {},
      topics: [],
    }
  }

  componentDidMount() {
    this.api.user.show(this.state.id).then((resp) => {
      if (resp.error) {
        return
      }
      let user = resp.data;
      user.created_at = moment(user.created_at).format('l');
      user.biography = user.biography.slice(0, 256);
      this.setState({user: user}, () => {
        this.api.user.topics(this.state.id).then((resp) => {
          if (resp.error) {
            return
          }
          this.setState({topics: resp.data});
        });
      });
    });
  }

  render() {
    const i18n = window.i18n;
    let state = this.state;
    const topics = state.topics.map((topic) => {
      return (
        <TopicItem topic={topic} key={topic.topic_id} profile={true}/>
      )
    });

    const profile = (
      <div className={style.profile}>
        <img src={state.user.avatar_url} className={style.avatar} alt={state.user.nickname} />
        <div className={style.name}>
          {state.user.nickname}
        </div>
        <div className={style.created}>
          Joined {state.user.created_at}
        </div>
        <div className={style.biography}>
          {state.user.biography}
        </div>
      </div>
    );

    return (
      <div className='container'>
        <aside className='column aside'>
          {profile}
        </aside>
        <main className='column main'>
          <div className={style.title}>
              {i18n.t('user.topics')}
          </div>
          <div className={style.topics}>
            <ul>
              {topics}
            </ul>
          </div>
        </main>
      </div>
    )
  }
}

export default Show;
