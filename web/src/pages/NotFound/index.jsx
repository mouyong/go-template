import { Button, Typography } from '@douyinfe/semi-ui';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

const { Title, Text } = Typography;

export default function NotFound() {
  const { t } = useTranslation();
  const navigate = useNavigate();

  return (
    <div className='min-h-screen flex flex-col items-center justify-center'>
      <Title heading={1} className='text-6xl mb-4'>
        404
      </Title>
      <Text size='large' className='mb-8'>
        {t('notFound.message')}
      </Text>
      <Button type='primary' onClick={() => navigate('/')}>
        {t('notFound.backHome')}
      </Button>
    </div>
  );
}
