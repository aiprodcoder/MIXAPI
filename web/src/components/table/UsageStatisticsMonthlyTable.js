import React, { useEffect, useState } from 'react';
import {
  API,
  showError,
  showSuccess,
  renderQuota,
  isAdmin,
} from '../../helpers';
import { ITEMS_PER_PAGE } from '../../constants';
import {
  Button,
  Card,
  Form,
  Modal,
  Space,
  Table,
  Tag,
  Typography,
  DatePicker,
  Select,
  Row,
  Col,
} from '@douyinfe/semi-ui';
import { IconSearch, IconRefresh, IconCalendar } from '@douyinfe/semi-icons';
import { useTranslation } from 'react-i18next';
import { useTableCompactMode } from '../../hooks/useTableCompactMode';

const { Text, Title } = Typography;

const UsageStatisticsMonthlyTable = () => {
  const { t } = useTranslation();

  const columns = [
    {
      title: t('月份'),
      dataIndex: 'date',
      key: 'date',
      render: (text) => {
        // 只显示年月部分 YYYY-MM
        const month = text.substring(0, 7);
        return <Text strong>{month}</Text>;
      },
      width: 120,
    },
    {
      title: t('令牌名称'),
      dataIndex: 'token_name',
      key: 'token_name',
      render: (text, record) => (
        <div>
          <Text>{text || t('未知令牌')}</Text>
          <br />
          <Text type="tertiary" size="small">ID: {record.token_id}</Text>
        </div>
      ),
      width: 150,
    },
    {
      title: t('模型名称'),
      dataIndex: 'model_name',
      key: 'model_name',
      render: (text) => (
        <Tag color="blue" shape="circle">
          {text}
        </Tag>
      ),
      width: 150,
    },
    {
      title: t('请求统计'),
      key: 'requests',
      render: (text, record) => (
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-2">
            <Text size="small" type="tertiary">{t('总数')}:</Text>
            <Text strong>{record.total_requests}</Text>
          </div>
          <div className="flex items-center gap-2">
            <Text size="small" type="tertiary">{t('成功')}:</Text>
            <Text type="success">{record.successful_requests}</Text>
          </div>
          <div className="flex items-center gap-2">
            <Text size="small" type="tertiary">{t('失败')}:</Text>
            <Text type="danger">{record.failed_requests}</Text>
          </div>
        </div>
      ),
      width: 120,
    },
    {
      title: t('成功率'),
      key: 'success_rate',
      render: (text, record) => {
        const rate = record.total_requests > 0 
          ? ((record.successful_requests / record.total_requests) * 100).toFixed(1)
          : '0.0';
        const color = parseFloat(rate) >= 95 ? 'green' : parseFloat(rate) >= 80 ? 'orange' : 'red';
        return (
          <Tag color={color} shape="circle">
            {rate}%
          </Tag>
        );
      },
      width: 100,
    },
    {
      title: t('Token统计'),
      key: 'tokens',
      render: (text, record) => (
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-2">
            <Text size="small" type="tertiary">{t('总计')}:</Text>
            <Text strong>{record.total_tokens.toLocaleString()}</Text>
          </div>
          <div className="flex items-center gap-2">
            <Text size="small" type="tertiary">{t('输入')}:</Text>
            <Text>{record.prompt_tokens.toLocaleString()}</Text>
          </div>
          <div className="flex items-center gap-2">
            <Text size="small" type="tertiary">{t('输出')}:</Text>
            <Text>{record.completion_tokens.toLocaleString()}</Text>
          </div>
        </div>
      ),
      width: 150,
    },
    {
      title: t('额度消耗'),
      dataIndex: 'total_quota',
      key: 'total_quota',
      render: (text) => (
        <Text strong type="warning">
          {renderQuota(text)}
        </Text>
      ),
      width: 120,
    },
    {
      title: t('更新时间'),
      dataIndex: 'updated_time',
      key: 'updated_time',
      render: (text) => {
        const date = new Date(text * 1000);
        return (
          <Text size="small" type="tertiary">
            {date.toLocaleString()}
          </Text>
        );
      },
      width: 150,
    },
  ];

  const [pageSize, setPageSize] = useState(ITEMS_PER_PAGE);
  const [statistics, setStatistics] = useState([]);
  const [loading, setLoading] = useState(true);
  const [activePage, setActivePage] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [compactMode, setCompactMode] = useTableCompactMode('usage_statistics_monthly');
  const [summary, setSummary] = useState(null);
  const [tokens, setTokens] = useState([]);

  // Form 初始值
  const getInitValues = () => {
    const now = new Date();
    const sixMonthsAgo = new Date(now.getFullYear(), now.getMonth() - 6, 1);
    return {
      startDate: sixMonthsAgo,
      endDate: now,
      tokenId: '',
      modelName: '',
    };
  };

  // Form API 引用
  const [formApi, setFormApi] = useState(null);

  // 获取表单值的辅助函数
  const getFormValues = () => {
    const formValues = formApi ? formApi.getValues() : {};
    const startDate = formValues.startDate ? 
      new Date(formValues.startDate).toISOString().split('T')[0].substring(0, 7) : '';
    const endDate = formValues.endDate ? 
      new Date(formValues.endDate).toISOString().split('T')[0].substring(0, 7) : '';
    
    return {
      start_date: startDate,
      end_date: endDate,
      token_id: formValues.tokenId || '',
      model_name: formValues.modelName || '',
    };
  };

  // 将后端返回的数据写入状态
  const syncPageData = (payload) => {
    setStatistics(payload.items || []);
    setTotalCount(payload.total || 0);
    setActivePage(payload.page || 1);
    setPageSize(payload.page_size || pageSize);
    setSummary(payload.summary || null);
  };

  const loadStatistics = async (page = 1, size = pageSize, searchParams = {}) => {
    setLoading(true);
    try {
      const params = new URLSearchParams({
        p: page.toString(),
        size: size.toString(),
        ...searchParams,
      });
      
      // 根据用户权限选择API端点
      const apiEndpoint = isAdmin() 
        ? `/api/usage_statistics_monthly/?${params}`
        : `/api/usage_statistics_monthly/self?${params}`;
      
      console.log('Loading monthly statistics with params:', params.toString());
      console.log('API endpoint:', apiEndpoint);
      const res = await API.get(apiEndpoint);
      console.log('API response:', res.data);
      const { success, message, data } = res.data;
      if (success) {
        syncPageData(data);
      } else {
        showError(message);
      }
    } catch (error) {
      console.error('Failed to load monthly statistics:', error);
      showError(t('加载数据失败'));
    }
    setLoading(false);
  };

  const loadTokens = async () => {
    try {
      const res = await API.get('/api/token/');
      const { success, data } = res.data;
      if (success) {
        const tokenOptions = data.items.map(token => ({
          label: token.name,
          value: token.id,
        }));
        setTokens(tokenOptions);
      }
    } catch (error) {
      console.error('Failed to load tokens:', error);
      // 即使加载令牌失败，也不影响统计数据的显示
    }
  };

  const refresh = async (page = activePage) => {
    const searchParams = getFormValues();
    await loadStatistics(page, pageSize, searchParams);
  };

  const searchStatistics = async () => {
    const searchParams = getFormValues();
    await loadStatistics(1, pageSize, searchParams);
  };

  const resetSearch = () => {
    if (formApi) {
      formApi.setValues(getInitValues());
      loadStatistics(1, pageSize);
    }
  };

  useEffect(() => {
    console.log('Component mounted, user is admin:', isAdmin());
    loadStatistics(1);
    loadTokens();
  }, [pageSize]);

  return (
    <div className="space-y-4">
      {/* 调试信息 */}
      {process.env.NODE_ENV === 'development' && (
        <Card className="!rounded-2xl shadow-sm border-0 bg-yellow-50">
          <Text size="small">
            调试信息: 用户权限={isAdmin() ? '管理员' : '普通用户'}, 
            数据数量={statistics.length}, 
            总数={totalCount}, 
            加载状态={loading ? '加载中' : '完成'}
          </Text>
        </Card>
      )}
      {/* 统计摘要卡片 */}
      {summary && (
        <Card className="!rounded-2xl shadow-sm border-0">
          <div className="flex items-center mb-4">
            <IconCalendar className="mr-2 text-blue-500" size={20} />
            <Title heading={5} className="m-0">
              {t('统计摘要')}
            </Title>
          </div>
          <Row gutter={16}>
            <Col span={6}>
              <div className="statistic-card">
                <Text type="tertiary" size="small">{t('总请求数')}</Text>
                <Title heading={3} style={{ color: '#3f6600', margin: '4px 0 0 0' }}>
                  {summary.total_requests}
                </Title>
              </div>
            </Col>
            <Col span={6}>
              <div className="statistic-card">
                <Text type="tertiary" size="small">{t('成功请求数')}</Text>
                <Title heading={3} style={{ color: '#52c41a', margin: '4px 0 0 0' }}>
                  {summary.successful_requests}
                </Title>
              </div>
            </Col>
            <Col span={6}>
              <div className="statistic-card">
                <Text type="tertiary" size="small">{t('成功率')}</Text>
                <Title heading={3} style={{ 
                  color: summary.success_rate >= 95 ? '#52c41a' : '#fa8c16', 
                  margin: '4px 0 0 0' 
                }}>
                  {summary.success_rate.toFixed(1)}%
                </Title>
              </div>
            </Col>
            <Col span={6}>
              <div className="statistic-card">
                <Text type="tertiary" size="small">{t('总额度消耗')}</Text>
                <Title heading={3} style={{ color: '#1890ff', margin: '4px 0 0 0' }}>
                  {renderQuota(summary.total_quota)}
                </Title>
              </div>
            </Col>
          </Row>
        </Card>
      )}

      {/* 搜索区域 */}
      <Card className="!rounded-2xl shadow-sm border-0">
        <Form
          getFormApi={(api) => setFormApi(api)}
          initValues={getInitValues()}
          onSubmit={searchStatistics}
          layout="horizontal"
          className="search-form"
        >
          <Row gutter={16}>
            <Col span={6}>
              <Form.DatePicker
                field="startDate"
                label={t('开始月份')}
                style={{ width: '100%' }}
                placeholder={t('请选择开始月份')}
                picker="month"
              />
            </Col>
            <Col span={6}>
              <Form.DatePicker
                field="endDate"
                label={t('结束月份')}
                style={{ width: '100%' }}
                placeholder={t('请选择结束月份')}
                picker="month"
              />
            </Col>
            <Col span={6}>
              <Form.Select
                field="tokenId"
                label={t('令牌')}
                style={{ width: '100%' }}
                placeholder={t('请选择令牌')}
                optionList={tokens}
                showClear
                filter
              />
            </Col>
            <Col span={6}>
              <Form.Input
                field="modelName"
                label={t('模型名称')}
                style={{ width: '100%' }}
                placeholder={t('请输入模型名称')}
                showClear
              />
            </Col>
          </Row>
          <div className="flex justify-end mt-4 gap-2">
            <Button
              type="primary"
              htmlType="submit"
              icon={<IconSearch />}
              loading={loading}
            >
              {t('搜索')}
            </Button>
            <Button
              icon={<IconRefresh />}
              onClick={resetSearch}
            >
              {t('重置')}
            </Button>
          </div>
        </Form>
      </Card>

      {/* 数据表格 */}
      <Card className="!rounded-2xl shadow-sm border-0">
        <div className="flex justify-between items-center mb-4">
          <div className="flex items-center">
            <Title heading={5} className="m-0">
              {t('用量月统计')}
            </Title>
            <Text type="tertiary" className="ml-2">
              ({t('按月份、令牌、模型分组汇总')})
            </Text>
          </div>
          <div className="flex items-center gap-2">
            <Button
              icon={<IconRefresh />}
              onClick={() => refresh()}
              loading={loading}
            >
              {t('刷新')}
            </Button>
          </div>
        </div>

        <Table
          columns={columns}
          dataSource={statistics}
          empty={
            statistics.length === 0 && !loading ? (
              <div className="text-center py-8">
                <Text type="tertiary">
                  {totalCount === 0 ? t('暂无统计数据，请先发起API请求或调整筛选条件') : t('暂无数据')}
                </Text>
              </div>
            ) : undefined
          }
          pagination={{
            currentPage: activePage,
            pageSize: pageSize,
            total: totalCount,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total, range) => 
              t('第 {{start}} - {{end}} 条，共 {{total}} 条', {
                start: range[0],
                end: range[1], 
                total: total
              }),
            onPageChange: (page) => {
              setActivePage(page);
              refresh(page);
            },
            onPageSizeChange: (size) => {
              setPageSize(size);
              setActivePage(1);
              loadStatistics(1, size, getFormValues());
            },
          }}
          loading={loading}
          size={compactMode ? 'small' : 'default'}
          rowKey="id"
          scroll={{ x: 1200 }}
        />
      </Card>
    </div>
  );
};

export default UsageStatisticsMonthlyTable;