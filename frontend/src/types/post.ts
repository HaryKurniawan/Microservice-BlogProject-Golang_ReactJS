export interface Post {
  id: number;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
  author_id?: number;
  author_name?: string;
  author_email?: string;
}

export interface CreatePostPayload {
  title: string;
  content: string;
}

export interface UpdatePostPayload {
  title: string;
  content: string;
}
