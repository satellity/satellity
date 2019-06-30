import style from './show.scss';
import React, { Component } from 'react';
import { Helmet } from 'react-helmet';
import API from '../api/index.js';
import Config from '../components/config.js';

class Show extends Component {
  constructor(props) {
    super(props);
    this.api = new API();

    let id = this.props.match.params.id;
    this.state = {
      group_id: id,
      name: '',
      description: '',
      user: {},
      loading: true
    }
  }

  componentDidMount() {
    this.api.group.show(this.state.group_id).then((data) => {
      data.loading = false;
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

    return (
      <div className='container'>
        {!state.loading && seoView}
        <main className='column main'>
          <div className={style.group}>
            <h1>{state.name}</h1>
            <div>
              {state.description}
            </div>
          </div>
          <div>
            <h4>{i18n.t('group.members')}</h4>
          </div>
          <div>
            <h4>{i18n.t('group.lastest_discussions')}</h4>
          </div>
        </main>
        <aside className='column aside'>
        </aside>
      </div>
    )
  }
}

export default Show;
