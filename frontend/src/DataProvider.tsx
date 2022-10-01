import React, {Component} from "react";
import {ArticleModel, ArticleResponse} from "./articles/ArticleModel";
import axios, {AxiosError, AxiosResponse} from "axios";
import {RepoModel, RepoResponse} from "./repos/RepoModel";
import {ContactResponse, ContactData} from "./contact/ContactModel";
import {FormResponse} from "./contact/ContactComponent";
import {getAnalyticsParams, logAnalyticsEvent} from "./analytics/firebase";
import {AnalyticsEvent} from "./analytics/events";
import {AnalyticsKey} from "./analytics/keys";

export type DataValue = {
  articles: ArticleModel[],
  repos: RepoModel[],
  fetchArticles: () => void,
  fetchRepos: () => void,
  postEmail: (data: ContactData) => Promise<FormResponse>
}

const DataContext = React.createContext<DataValue>({} as DataValue);

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
        console.error(e)
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
        logNetworkError(e, JSON.stringify(data.redact()))
        return DataProvider.generateContactResponse()
      })
  }

  private static generateContactResponse(success?: boolean): FormResponse {
    if (success === true) {
      return {message: "Your message has been received. I should revert within 24 hours.", success: true}
    }
    return {message: "Something went wrong. Please, try again", success: false}
  }

  render() {
    return (
      <DataContext.Provider
        value={{
          ...this.state,
          fetchArticles: this.fetchArticles,
          fetchRepos: this.fetchRepos,
          postEmail: this.postEmail,
        }}>
        {this.props.children}
      </DataContext.Provider>
    )
  }
}

function logNetworkError(e: AxiosError, data?: string) {
  let params = getAnalyticsParams()
  params.set(AnalyticsKey.url, `${e.config?.baseURL}${e.config?.url}`)
  params.set(AnalyticsKey.method, e.config?.method)
  params.set(AnalyticsKey.message, e.message)

  if (data) params.set(AnalyticsKey.data, data)
  if (e.response) {
    params.set(AnalyticsKey.statusCode, e.response.status)
    params.set(AnalyticsKey.rawResponse, JSON.stringify(e.response.data))
  }

  logAnalyticsEvent(AnalyticsEvent.apiFailure, params)
}


export {DataContext}
export default DataProvider
