import { useState } from 'react';
import { Form, Button } from '@douyinfe/semi-ui';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { toast } from 'react-toastify';
import { register } from '@/helpers/auth';

export default function RegisterForm() {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (values) => {
    setLoading(true);
    try {
      await register(values);
      toast.success(t('register.success'));
      navigate('/');
    } catch (error) {
      toast.error(error.message || t('register.failed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <Form onSubmit={handleSubmit} className='w-full max-w-md'>
      <Form.Input
        field='username'
        label={t('register.username')}
        placeholder={t('register.usernamePlaceholder')}
        rules={[{ required: true, message: t('register.usernameRequired') }]}
      />
      <Form.Input
        field='email'
        label={t('register.email')}
        placeholder={t('register.emailPlaceholder')}
        rules={[
          { required: true, message: t('register.emailRequired') },
          { type: 'email', message: t('register.emailInvalid') },
        ]}
      />
      <Form.Input
        field='password'
        label={t('register.password')}
        mode='password'
        placeholder={t('register.passwordPlaceholder')}
        rules={[
          { required: true, message: t('register.passwordRequired') },
          { min: 6, message: t('register.passwordMin') },
        ]}
      />
      <Button type='primary' htmlType='submit' loading={loading} block className='mt-4'>
        {t('register.submit')}
      </Button>
    </Form>
  );
}
