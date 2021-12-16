export interface ArticleAPI {
  data:    Article[];
  success: boolean;
}

export interface Article {
  title:      string;
  thumbnail:  string;
  url:        string;
  posted_at:  Date;
  updated_at: Date;
  excerpt:    string;
}
