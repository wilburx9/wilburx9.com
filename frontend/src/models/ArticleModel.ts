export interface ArticleResponse {
  data:    ArticleModel[];
  success: boolean;
}

export interface ArticleModel {
  title:      string;
  thumbnail:  string;
  url:        string;
  posted_at:  Date;
  updated_at: Date;
  excerpt:    string;
}
