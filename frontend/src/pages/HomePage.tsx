import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { notification } from 'antd';
import type { Post } from '../types/post';
import * as postApi from '../api/postApi';
import Header from '../components/Header';
import LoginPrompt from '../components/LoginPrompt';
import PostForm from '../components/PostForm';
import PostList from '../components/PostList';

export default function HomePage() {
  const navigate = useNavigate();
  const [posts, setPosts] = useState<Post[]>([]);
  const [editingPost, setEditingPost] = useState<Post | null>(null);
  const [listLoading, setListLoading] = useState(false);
  const [formLoading, setFormLoading] = useState(false);
  
  // State untuk Tab ("all" = Semua Post, "my" = Post Saya)
  const [activeTab, setActiveTab] = useState<'all' | 'my'>('all');

  // State Autentikasi
  const [user, setUser] = useState<any | null>(null);

  const fetchPosts = useCallback(async () => {
    setListLoading(true);
    try {
      setPosts(await postApi.getAllPosts());
    } catch {
      notification.error({
        title: 'Error',
        description: 'Gagal memuat posts. Pastikan API Gateway microservice sudah berjalan di port 8085.',
        placement: 'topRight',
      });
    } finally {
      setListLoading(false);
    }
  }, []);

  useEffect(() => {
    // Sinkronisasi status login langsung dari localStorage saat komponen dimuat
    const savedUser = localStorage.getItem('user');
    setUser(savedUser ? JSON.parse(savedUser) : null);
    
    fetchPosts();
  }, [fetchPosts]);

  const handleLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    setUser(null);
    setActiveTab('all'); // Balikkan ke tab all jika logout
    notification.success({
      title: 'Logout Berhasil',
      description: 'Anda telah keluar dari akun.',
      placement: 'topRight',
    });
  };

  const handleSubmit = async (payload: { title: string; content: string }) => {
    const currentToken = localStorage.getItem('token');
    if (!currentToken) {
      notification.warning({
        title: 'Perhatian',
        description: 'Silakan login terlebih dahulu untuk membuat postingan.',
        placement: 'topRight',
      });
      navigate('/auth');
      return;
    }

    setFormLoading(true);
    try {
      if (editingPost) {
        // 1. Kirim update ke server
        const updatedPost = await postApi.updatePost(editingPost.id, payload);
        
        // 2. Update state lokal
        setPosts(prevPosts =>
          prevPosts.map(p => (p.id === editingPost.id ? updatedPost : p))
        );

        notification.success({
          title: 'Berhasil',
          description: 'Post berhasil diupdate!',
          placement: 'topRight',
        });
        setEditingPost(null);
      } else {
        // 1. Kirim data baru ke server
        const newPost = await postApi.createPost(payload);
        
        // 2. Sisipkan langsung data baru ke awal array posts
        setPosts(prevPosts => [newPost, ...prevPosts]);

        notification.success({
          title: 'Berhasil',
          description: 'Post berhasil dibuat!',
          placement: 'topRight',
        });
      }
    } catch (err: any) {
      if (err.response?.status === 401) {
        notification.error({
          title: 'Sesi Berakhir',
          description: 'Sesi Anda telah habis. Silakan login kembali.',
          placement: 'topRight',
        });
        handleLogout();
        navigate('/auth');
      } else {
        notification.error({
          title: 'Error',
          description: 'Gagal menyimpan post. Silakan coba lagi.',
          placement: 'topRight',
        });
      }
    } finally {
      setFormLoading(false);
    }
  };

  const handleDelete = async (id: number) => {
    const currentToken = localStorage.getItem('token');
    if (!currentToken) {
      notification.warning({
        title: 'Perhatian',
        description: 'Silakan login terlebih dahulu untuk menghapus postingan.',
        placement: 'topRight',
      });
      navigate('/auth');
      return;
    }

    if (!confirm('Yakin ingin menghapus post ini?')) return;
    try {
      await postApi.deletePost(id);
      setPosts(prevPosts => prevPosts.filter(p => p.id !== id));
      notification.success({
        title: 'Berhasil',
        description: 'Post berhasil dihapus.',
        placement: 'topRight',
      });
    } catch (err: any) {
      if (err.response?.status === 401) {
        notification.error({
          title: 'Sesi Berakhir',
          description: 'Sesi Anda telah habis. Silakan login kembali.',
          placement: 'topRight',
        });
        handleLogout();
        navigate('/auth');
      } else {
        notification.error({
          title: 'Error',
          description: 'Gagal menghapus post.',
          placement: 'topRight',
        });
      }
    }
  };

  const handleEdit = (post: Post) => {
    const currentToken = localStorage.getItem('token');
    if (!currentToken) {
      navigate('/auth');
      return;
    }
    setEditingPost(post);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  // Filter posts berdasarkan tab yang aktif
  const displayedPosts = activeTab === 'all' 
    ? posts 
    : posts.filter(post => user && post.author_id === user.id);

  return (
    <div className="max-w-5xl mx-auto px-6 py-10 pb-20">
      {/* Header Utama */}
      <Header user={user} postsCount={posts.length} onLogout={handleLogout} />

      {/* Grid Konten */}
      <div className="grid grid-cols-1 lg:grid-cols-[360px_1fr] gap-8 items-start">
        {/* Kolom Kiri: Form (Jika Login) atau Promosi Login (Jika Belum) */}
        <section>
          {user ? (
            <PostForm
              editingPost={editingPost}
              onSubmit={handleSubmit}
              onCancel={() => setEditingPost(null)}
              loading={formLoading}
            />
          ) : (
            <LoginPrompt />
          )}
        </section>

        {/* Kolom Kanan: Daftar Artikel dengan Tabbing */}
        <section>
          <div className="flex flex-col gap-4 mb-6 border-b border-[var(--color-border)] pb-2">
            <div className="flex items-center justify-between flex-wrap gap-3">
              {/* Tab Selector Buttons */}
              <div className="flex gap-2">
                <button
                  onClick={() => setActiveTab('all')}
                  className={`
                    text-sm font-semibold pb-2 px-1 border-b-2 cursor-pointer transition-all duration-150
                    ${activeTab === 'all' 
                      ? 'border-[var(--color-accent)] text-[var(--color-accent)]' 
                      : 'border-transparent text-[var(--color-text-muted)] hover:text-[var(--color-text-base)]'}
                  `}
                >
                  Semua Post
                </button>
                
                {user && (
                  <button
                    onClick={() => setActiveTab('my')}
                    className={`
                      text-sm font-semibold pb-2 px-1 border-b-2 cursor-pointer transition-all duration-150
                      ${activeTab === 'my' 
                        ? 'border-[var(--color-accent)] text-[var(--color-accent)]' 
                        : 'border-transparent text-[var(--color-text-muted)] hover:text-[var(--color-text-base)]'}
                    `}
                  >
                    Post Saya
                  </button>
                )}
              </div>

              {/* Refresh Button */}
              <button
                onClick={fetchPosts}
                disabled={listLoading}
                className="
                  text-xs font-medium px-3 py-1.5 rounded-[var(--radius-btn)]
                  bg-[var(--color-bg-elevated)] text-[var(--color-text-muted)]
                  hover:text-[var(--color-text-base)] hover:bg-[var(--color-border)]
                  transition-all duration-150 cursor-pointer disabled:opacity-50
                "
              >
                {listLoading ? 'Memuat...' : '↻ Segarkan'}
              </button>
            </div>
          </div>

          <PostList
            posts={displayedPosts}
            loading={listLoading}
            onEdit={handleEdit}
            onDelete={handleDelete}
          />
        </section>
      </div>
    </div>
  );
}
