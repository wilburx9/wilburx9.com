import React, {Component} from "react";
import {ArticleModel, ArticleResponse} from "./models/ArticleModel";
import axios, {AxiosError, AxiosResponse} from "axios";
import {RepoModel, RepoResponse} from "./models/RepoModel";
import {ContactResponse, ContactData} from "./models/ContactModel";
import {FormResponse} from "./components/ContactComponent";
import {getAnalyticsParams, logAnalyticsEvent} from "./analytics/firebase";
import {AnalyticsEvent} from "./analytics/events";
import {AnalyticsKey} from "./analytics/keys";

export type DataValue = {
  articles: ArticleModel[],
  repos: RepoModel[],
  fetchArticles: () => void,
  fetchRepos: () => void,
  postEmail: (data: ContactData) => Promise<FormResponse>,
  hasData: () => boolean,
}

const DataContext = React.createContext<Partial<DataValue>>({});

const http = axios.create({
  baseURL: `${process.env.REACT_APP_DOMAIN}/api`,
  headers: {"Content-type": "application/json"}
})

type DataState = {
  articles: ArticleModel[],
  repos: RepoModel[],
}

export class DataProvider extends Component<any, DataState> {

  state: DataState = {articles: [], repos: []}

  componentDidMount() {
    this.fetchArticles()
    this.fetchRepos()
  }

  fetchArticles = () => {
    http
      .get<ArticleResponse>("/articles")
      .then(response => {
        this.setState({articles: response.data.data})
      })
      .catch(e => {
        logNetworkError(e)
      })
  }

  fetchRepos = () => {
    let params = new URLSearchParams([["size", "6"]])
    http
      .get<RepoResponse>("/repos", {params})
      .then(response => {
        this.setState({repos: response.data.data})
      })
      .catch(e => {
        logNetworkError(e)
      })
  }

  postEmail = async (data: ContactData): Promise<FormResponse> => {
    return http
      .post<ContactData, AxiosResponse<ContactResponse>>("/contact", data)
      .then(response => {
        return DataProvider.generateContactResponse(response.data.success)
      })
      .catch((e) => {
        logNetworkError(e)
        return DataProvider.generateContactResponse()
      })
  }

  private static generateContactResponse(success?: boolean): FormResponse {
    if (success === true) {
      return {message: "Your message has been received. I will reply as soon as I can.", success: true}
    }
    return {message: "Something went wrong. Please, try again", success: false}
  }

  hasData = (): boolean => {
    return this.state.repos.length > 0 || this.state.articles.length > 0
  }

  render() {
    return (
      <DataContext.Provider
        value={{
          ...this.state,
          fetchArticles: this.fetchArticles,
          fetchRepos: this.fetchRepos,
          postEmail: this.postEmail,
          hasData: this.hasData
        }}>
        {this.props.children}
      </DataContext.Provider>
    )
  }
}

function logNetworkError(e: AxiosError) {
  let params = getAnalyticsParams()
  params.set(AnalyticsKey.url, `${e.config.baseURL}${e.config.url}`)
  params.set(AnalyticsKey.method, e.config.method)
  params.set(AnalyticsKey.message, e.message)
  params.set(AnalyticsKey.data, e.config.data)

  if (e.response) {
    params.set(AnalyticsKey.statusCode,  e.response.status)
    params.set(AnalyticsKey.rawResponse, JSON.stringify(e.response.data))
  }

  logAnalyticsEvent(AnalyticsEvent.apiFailure, params)
}


export {DataContext}
export default DataProvider