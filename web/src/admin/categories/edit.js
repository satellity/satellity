import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from '../../api/index.js';
import LoadingView from '../../loading/loading.js';

class AdminCategoryEdit extends Component {
  constructor(props) {
    super(props);
    this.api = new API();
    this.id = props.match.params.id;
    this.handleSubmit = this.handleSubmit.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.state = {type: 'category', name: '', alias: '', position: '', description: '', submitting: false};
  }

  componentDidMount() {
    const self = this;
    this.api.category.show(self.id, (resp) => {
      self.setState(resp.data);
    });
  }

  handleChange(e) {
    const target = e.target;
    const name = target.name;
    this.setState({
      [name]: target.value
    });
  }

  handleSubmit(e) {
    e.preventDefault();
    this.setState({submitting: true});
    const history = this.props.history;
    this.api.category.update(this.id, this.state, () => {
      history.push('/admin/categories');
      this.setState({submitting: false});
    });
  }

  render() {
    return (
      <CategoryEdit onSubmit={this.handleSubmit} onChange={this.handleChange} state={this.state} />
    )
  }
}

const CategoryEdit = (props) => (
  <div>
    <h1 className='welcome'>
      Edit a category
    </h1>
    <div className='panel'>
      <form onSubmit={(e) => props.onSubmit(e)}>
        <div>
          <label name='name'>Name *</label>
          <input type='text' name='name' value={props.state.name} autoComplete='off' onChange={(e) => props.onChange(e)}/>
        </div>
        <div>
          <label name='alias'>Alias *</label>
          <input type='text' name='alias' value={props.state.alias} autoComplete='off' onChange={(e) => props.onChange(e)}/>
        </div>
        <div>
          <label name='position'>Position</label>
          <input type='number' name='position' value={props.state.position} autoComplete='off' onChange={(e) => props.onChange(e)}/>
        </div>
        <div>
          <label name='name'>Description *</label>
          <textarea type='text' name='description' value={props.state.description} onChange={(e) => props.onChange(e)} key='description'/>
        </div>
        <div className='action'>
          <button className='btn submit blue' disabled={props.state.submitting}>
            { props.state.submitting && <LoadingView style='sm-ring blank'/> }
            &nbsp;SUBMIT
          </button>
        </div>
      </form>
    </div>
  </div>
)

export default AdminCategoryEdit;
