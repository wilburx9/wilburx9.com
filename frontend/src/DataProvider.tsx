import React, {Component} from "react";
import {Article, ArticleAPI} from "./models/Article";
import axios from "axios";

export type DataValue = {
  articles: Article[],
  fetchArticles: () => void
}

const DataContext = React.createContext<Partial<DataValue>>({});

const http = axios.create({
  baseURL: `${process.env.REACT_APP_DOMAIN}/api`,
  headers: {"Content-type": "application/json"}
})

type DataState = {
  articles: Article[]
}

export class DataProvider extends Component<any, DataState> {

  state: DataState = {articles: []}

  componentDidMount() {
    this.fetchArticles()
  }

  fetchArticles = () => {
    http
      .get<ArticleAPI>("/articles")
      .then(response => {
        console.log("Success:: " + response.data)
        this.setState({articles: response.data.data})
      })
      .catch(ex => {
        console.log(ex)
      })
  }

  render() {
    return (
      <DataContext.Provider
        value={{
          ...this.state,
          fetchArticles: this.fetchArticles
        }}>
        {this.props.children}
      </DataContext.Provider>
    )
  }
}

const DataConsumer = DataContext.Consumer

export {DataConsumer, DataContext}
export default DataProvider