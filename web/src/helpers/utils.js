import dayjs from 'dayjs';
import relativeTime from 'dayjs/plugin/relativeTime';
import 'dayjs/locale/zh-cn';

dayjs.extend(relativeTime);
dayjs.locale('zh-cn');

export const formatDate = (date, format = 'YYYY-MM-DD HH:mm:ss') => {
  return dayjs(date).format(format);
};

export const formatTime = (timestamp, relative = false) => {
  if (!timestamp) return '';

  // 如果是秒级时间戳，转换为毫秒
  const ts = timestamp < 10000000000 ? timestamp * 1000 : timestamp;
  const target = dayjs(ts);

  // 如果不需要相对时间，直接返回格式化的日期
  if (!relative) {
    return target.format('YYYY-MM-DD HH:mm');
  }

  // 相对时间展示
  const now = dayjs();
  const diffInSeconds = now.diff(target, 'second');

  // 1分钟内显示"刚刚"
  if (diffInSeconds < 60) {
    return '刚刚';
  }

  // 1小时内显示"X分钟前"
  if (diffInSeconds < 3600) {
    const minutes = Math.floor(diffInSeconds / 60);
    return `${minutes}分钟前`;
  }

  // 24小时内显示"X小时前"
  if (diffInSeconds < 86400) {
    const hours = Math.floor(diffInSeconds / 3600);
    return `${hours}小时前`;
  }

  // 7天内显示"X天前"
  if (diffInSeconds < 604800) {
    const days = Math.floor(diffInSeconds / 86400);
    return `${days}天前`;
  }

  // 超过7天显示具体日期
  return target.format('YYYY-MM-DD HH:mm');
};

export const formatNumber = (num) => {
  return new Intl.NumberFormat().format(num);
};

export const formatCurrency = (amount, currency = 'CNY') => {
  return new Intl.NumberFormat('zh-CN', {
    style: 'currency',
    currency,
  }).format(amount);
};

export const truncate = (str, length = 50) => {
  if (!str) return '';
  return str.length > length ? str.substring(0, length) + '...' : str;
};

export const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

export const debounce = (func, wait) => {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
};

export const throttle = (func, limit) => {
  let inThrottle;
  return function executedFunction(...args) {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
};
