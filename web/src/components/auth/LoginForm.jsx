import { useState } from 'react';
import { Form, Button } from '@douyinfe/semi-ui';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { toast } from 'react-toastify';
import { login } from '@/helpers/auth';

export default function LoginForm() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      await login(values);
      toast.success(t('login.success'));
      navigate('/dashboard');
    } catch (error) {
      toast.error(error.message || t('login.failed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Form onSubmit={handleSubmit} className='w-full max-w-md'>
      <Form.Input
        field='username'
        label={t('login.username')}
        placeholder={t('login.usernamePlaceholder')}
        rules={[{ required: true, message: t('login.usernameRequired') }]}
      />
      <Form.Input
        field='password'
        label={t('login.password')}
        mode='password'
        placeholder={t('login.passwordPlaceholder')}
        rules={[{ required: true, message: t('login.passwordRequired') }]}
      />
      <Button type='primary' htmlType='submit' loading={loading} block className='mt-4'>
        {t('login.submit')}
      </Button>
    </Form>
  );
}
