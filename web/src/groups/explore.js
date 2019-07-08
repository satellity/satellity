import style from './explore.scss';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../api/index.js';
import Item from './item.js';

class Explore extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.state = {groups: []};
  }

  componentDidMount() {
    this.api.group.index().then((data) => {
      this.setState({groups: data});
    });
  }

  render() {
    const groups = this.state.groups.map((group) => {
      return (
        <div key={group.group_id} className={style.item}>
          <Item group={group} />
        </div>
      )
    });

    return (
      <div className='wrapper container'>
        <div className={style.explore}>
          <h1>
            {i18n.t('group.explore')}
            <Link to='/groups/new' className={style.navi}>
              <FontAwesomeIcon icon={['fa', 'plus']} />
            </Link>
          </h1>
          <div className={style.list}>
            {groups}
          </div>
        </div>
      </div>
    )
  }
}

export default Explore;
