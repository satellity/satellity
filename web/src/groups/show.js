import style from './show.scss';
import React, { Component } from 'react';
import { Helmet } from 'react-helmet';
import {Link} from 'react-router-dom';
import API from '../api/index.js';
import Config from '../components/config.js';
import LoadingView from '../loading/loading.js';

class Show extends Component {
  constructor(props) {
    super(props);
    this.api = new API();

    let id = this.props.match.params.id;
    this.state = {
      group_id: id,
      name: '',
      description: '',
      owner: false,
      user: {},
      loading: true
    }
  }

  componentDidMount() {
    let me = this.api.user.me();
    this.api.group.show(this.state.group_id).then((data) => {
      data.loading = false;
      if (me && me.user_id == data.user.user_id) {
        data.owner = true;
      }
      this.setState(data);
    });
  }

  render() {
    const state = this.state;

    const seoView = (
      <Helmet>
        <title>{state.name} - {state.user.nickname} - {Config.Name}</title>
        <meta name='description' content={state.description.slice(0, 256)} />
      </Helmet>
    )

    const loadingView = (
      <div className={style.loading}>
        <LoadingView style='md-ring'/>
      </div>
    )

    const showView = (
      <div className={style.group}>
        <div className={style.head}>
          <div className={style.title}>
            <h1 className={style.name}>{state.name}</h1>
            <div className={style.nickname}>{state.user.nickname}</div>
          </div>
          <img src={state.user.avatar_url} className={style.avatar} />
        </div>
        <div>
          {state.description}
        </div>
      </div>
    )

    return (
      <div className='container'>
        {!state.loading && seoView}
        <main className='column main'>
          {state.loading && loadingView}
          {!state.loading && showView}
        </main>
        <aside className='column aside'>
          <div>
            <Link to={`/groups/${state.group_id}/members`}>
              {i18n.t('group.navi.members', {count: state.users_count})}
            </Link>
          </div>
          <div>
            <Link to={`/groups/${state.group_id}/messages`}>
              {i18n.t('group.navi.messages')}
            </Link>
          </div>
        </aside>
      </div>
    )
  }
}

export default Show;
