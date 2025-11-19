import { Card, Typography } from '@douyinfe/semi-ui';
import { useTranslation } from 'react-i18next';

const { Title, Text } = Typography;

export default function Dashboard() {
  const { t } = useTranslation();

  return (
    <div>
      <Title heading={2} className='mb-6'>
        {t('dashboard.title')}
      </Title>
      <div className='grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4'>
        <Card title={t('dashboard.users')} className='text-center'>
          <Text size='large' strong>
            1,234
          </Text>
        </Card>
        <Card title={t('dashboard.orders')} className='text-center'>
          <Text size='large' strong>
            567
          </Text>
        </Card>
        <Card title={t('dashboard.revenue')} className='text-center'>
          <Text size='large' strong>
            $12,345
          </Text>
        </Card>
        <Card title={t('dashboard.growth')} className='text-center'>
          <Text size='large' strong>
            +23%
          </Text>
        </Card>
      </div>
    </div>
  );
}
