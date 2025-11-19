import { Button } from '@douyinfe/semi-ui';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

export default function Home() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <div className='min-h-screen flex flex-col items-center justify-center bg-gradient-to-b from-blue-50 to-white'>
      <div className='text-center'>
        <h1 className='text-4xl font-bold text-gray-900 mb-4'>{t('home.title')}</h1>
        <p className='text-lg text-gray-600 mb-8'>{t('home.subtitle')}</p>
        <div className='space-x-4'>
          <Button type='primary' size='large' onClick={() => navigate('/dashboard')}>
            {t('home.getStarted')}
          </Button>
          <Button size='large' onClick={() => navigate('/about')}>
            {t('home.learnMore')}
          </Button>
        </div>
      </div>
    </div>
  );
}
