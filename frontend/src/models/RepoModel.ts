export interface RepoResponse {
  data: RepoModel[];
  success: boolean;
}

export interface RepoModel {
  name: string;
  stars: number;
  forks: number;
  url: string;
  description?: null | string;
  createdAt: Date;
  updatedAt: Date;
  license?: null | string;
  languages: Language[];
}

export interface Language {
  name: string;
  color: string;
}
