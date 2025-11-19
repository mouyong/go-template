import { lazy, Suspense } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import PageLayout from './components/layout/PageLayout';
import Loading from './components/common/Loading';

// 懒加载页面
const Home = lazy(() => import('./pages/Home'));
const Dashboard = lazy(() => import('./pages/Dashboard'));
const NotFound = lazy(() => import('./pages/NotFound'));

function App() {
  return (
    <Suspense fallback={<Loading />}>
      <Routes>
        {/* 公开路由 */}
        <Route path='/' element={<Home />} />

        {/* 需要布局的路由 */}
        <Route element={<PageLayout />}>
          <Route path='/dashboard' element={<Dashboard />} />
        </Route>

        {/* 404 */}
        <Route path='/404' element={<NotFound />} />
        <Route path='*' element={<Navigate to='/404' replace />} />
      </Routes>
    </Suspense>
  );
}

export default App;
