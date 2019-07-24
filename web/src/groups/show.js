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
    let user = this.api.user.local();
    this.api.group.show(this.state.group_id).then((data) => {
      data.loading = false;
      if (user && user.user_id == data.user.user_id) {
        data.owner = true;
      }
      this.setState(data);
    });
  }

  render() {
    let state = this.state;

    let seoView = (
      <Helmet>
        <title>{state.name} - {state.user.nickname} - {Config.Name}</title>
        <meta name='description' content={state.description.slice(0, 256)} />
      </Helmet>
    )

    let loadingView = (
      <div className={style.loading}>
        <LoadingView style='md-ring'/>
      </div>
    )

    let showView = (
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

    let sideView = (
      <div>
        <div className={style.navi}>
          <Link to={`/groups/${state.group_id}/members`}>
            {i18n.t('group.navi.members', {count: state.users_count})}
          </Link>
        </div>
        <div className={style.navi}>
          <Link to={`/groups/${state.group_id}/messages`}>
            {i18n.t('group.navi.messages')}
          </Link>
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
          {this.api.user.loggedIn() && sideView}
        </aside>
      </div>
    )
  }
}

export default Show;
