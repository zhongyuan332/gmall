import React, { Component }      from 'react';
import { connect }               from 'react-redux';

import {
    Button,
    Row, 
    Col,
    Form,
    Icon,
    Input,
    InputNumber,
    message,
    Modal,
    Popover,
    Select,
    TreeSelect,
    Tabs,
    Upload
} from 'antd';

import {
    requestProduct,
    requestSaveProduct,
    requestSaveProperty,
    requestSavePropertyValue,
    requestUpdateInventory,
    requestSaveInventory,
    requestUpdateHasProperty,
    requestUpdateTotalInventoryTemp,
    requestUpdateTotalInventory
} from '../../actions/product';

import requestCategoryList       from '../../actions/category/requestCategoryList';
import Software                  from '../Software';
import utils                     from '../../utils';
import analyze                   from '../../../sdk/analyze';
import '../../../../styles/admin/product/editProduct.css';

/*
 * 管理后台，新建或编辑商品
 */
class EditProduct extends Component {
    constructor(props) {
        super(props);
        this.onNameBlur               = this.onNameBlur.bind(this);
        this.onCategoriesChange       = this.onCategoriesChange.bind(this);
        this.onOriginalPriceBlur      = this.onOriginalPriceBlur.bind(this);
        this.onPriceBlur              = this.onPriceBlur.bind(this);
        this.onRemarkBlur             = this.onRemarkBlur.bind(this);
        this.onStatusChange           = this.onStatusChange.bind(this);
        this.onImageChange            = this.onImageChange.bind(this);
        this.onImageListChange        = this.onImageListChange.bind(this);
        this.onPreview                = this.onPreview.bind(this);
        this.onCancelPreview          = this.onCancelPreview.bind(this);
        this.onContentTypeChange      = this.onContentTypeChange.bind(this);
        this.onAddContent             = this.onAddContent.bind(this);
        this.onContentTextChange      = this.onContentTextChange.bind(this);
        this.onContentImageChange     = this.onContentImageChange.bind(this);
        this.onSubmit                 = this.onSubmit.bind(this);
        this.onPropValueVisibleChange = this.onPropValueVisibleChange.bind(this);
        this.onPropVisibleChange      = this.onPropVisibleChange.bind(this);
        this.onPropInput              = this.onPropInput.bind(this);
        this.addProp                  = this.addProp.bind(this);
        this.cancelAddProp            = this.cancelAddProp.bind(this);
        this.onInventoryChange        = this.onInventoryChange.bind(this);
        this.onSaveInventory          = this.onSaveInventory.bind(this);
        this.onTabChange              = this.onTabChange.bind(this);
        this.onHasPropertyChange      = this.onHasPropertyChange.bind(this);
        this.onTotalInventoryChange   = this.onTotalInventoryChange.bind(this);
        this.onSaveTotalInventory     = this.onSaveTotalInventory.bind(this);

        this.state = {
            activeTabKey        : '1',
            productId           : this.props.routeParams.id,
            categories          : [], //产品所属的分类
            name                : '',
            originalPrice       : 0,
            price               : 0,
            remark              : '',
            status              : '3', //等待上架
            imageID             : '',
            imageData           : '',
            previewVisible      : false,
            previewImage        : '',
            imageIDs            : '[]',
            imageList           : [],
            contentType         : 'image',
            contents            : [],
            properties          : [],
            inventories         : [],
            totalInventory      : 0,
            propValueVisibleMap : {},
            propValueTemp       : '',
            propPopupVisible    : false,
            propTemp            : '',
            hasProperty         : false,
            hasPropertyValue    : false,
            isLoading           : true
        };
    }
    componentDidMount() {
        analyze.pv();
        const { dispatch } = this.props;
        if (this.state.productId) {
            dispatch(requestProduct(this.state.productId));
        }
        dispatch(requestCategoryList());
    }
    componentWillReceiveProps(nextProps) {
        var self          = this;
        var product       = nextProps.data.product;
        var allCategories = nextProps.data.categories;

        function onDataReady(data) {
            console.log("onDataReady");
            var product    = data.product;
            var properties = product && product.properties || [];
            var propValueVisibleMap = {};
            var inventories = product && product.inventories || [];
            var hasPropertyValue = false;
            for (var i = 0; i < properties.length; i++) {
                propValueVisibleMap[properties[i].id] = false;
                if (properties[i].values && properties[i].values.length) {
                    hasPropertyValue = true;
                }
            }
            self.setState({
                productId           : product && product.id || '',
                categories          : data.categories || [],
                name                : product && product.name || '',
                contents            : product && JSON.parse(product.detail) || [],
                originalPrice       : product && product.originalPrice,
                price               : product && product.price,
                remark              : product && product.remark || '',
                status              : (product && product.status + '') || '3',
                imageID             : product && product.imageID || '',
                imageData           : data.imageURL || '',
                imageIDs            : product && product.imageIDs || '[]',
                imageList           : data.imageList || [],
                properties          : properties,
                inventories         : inventories,
                totalInventory      : product && product.totalInventory || 0,
                propValueVisibleMap : propValueVisibleMap,
                hasProperty         : product ? !!product.hasProperty : false,
                hasPropertyValue    : hasPropertyValue,
                isLoading           : false
            });
        }
        if (allCategories && allCategories.length > 0) {
            if (this.state.productId) {
                if (product) {
                    var categories = [];
                    for (var i = 0; i < product.categories.length; i++) {
                        var parentId = product.categories[i].parentId;
                        var id       = product.categories[i].id;
                        categories.push(utils.parseTreeNodeKey(allCategories, id));
                    }

                    var imageList  = [];
                    var pImageList = product.images || [];
                    for (var i = 0; i < pImageList.length; i++) {
                        imageList.push({
                            uid    : pImageList[i].id,
                            name   : pImageList[i].orignalTitle,
                            status : 'done',
                            url    : pImageList[i].url
                        });
                    }

                    var imageURL = product.image && product.image.url || '';
                    onDataReady({
                        product    : product,
                        imageURL   : imageURL,
                        imageList  : imageList,
                        categories : categories
                    });
                }
            } else {
                onDataReady({
                    product  : null,
                    imageURL : ''
                });
            }
        }
    }
    onNameBlur(event) {
        this.setState({ name: event.target.value });
    }
    onCategoriesChange(value) {
        this.setState({ categories: value });
    }
    onOriginalPriceBlur(event) {
        var price = event.target.value;
        price     = Number(price);
        this.setState({ originalPrice: price });
    }
    onPriceBlur(event) {
        var price = event.target.value;
        price     = Number(price);
        this.setState({ price: price });
    }
    onRemarkBlur(event) {
        this.setState({ remark: event.target.value });
    }
    onStatusChange(status) {
        this.setState({ status: status });
    }
    onBeforeUpload(file) {
        var isImage = file.type === 'image/jpeg' || file.type === 'image/png';
        if (!isImage) {
            message.error('只支持jpg或png格式的图片');
            return false;
        }
        const isLt2M = file.size / 1024 / 1024 < 2;
        if (!isLt2M) {
            message.error('图片大小要小于2M');
        }
        return isImage && isLt2M;
    }
    onImageChange(info) {
        var self = this;
        if (info.file.status === 'done') {
            self.setState({
                imageID  : info.file.response.data.id
            });
            (function(originFileObj, callback) {
                var reader = new FileReader();
                reader.addEventListener('load', function() {
                    callback(reader.result);
                });
                reader.readAsDataURL(originFileObj);
            }(info.file.originFileObj, function(imageData) {
                self.setState({ imageData })
            }));
        }
    }
    onCancelPreview() {
        this.setState({ previewVisible: false }); 
    }
    onPreview(file) {
        this.setState({
            previewImage   : file.url || file.thumbUrl,
            previewVisible : true
        });
    }
    onImageListChange(data) {
        var fileList = data.fileList || [];
        var imageIDs = [];
        for (var i = 0; i < fileList.length; i++) {
            if (fileList[i].response) {
                imageIDs.push(fileList[i].response.data.id);
            } else if (fileList[i].status == 'done') {
                imageIDs.push(fileList[i].uid);
            }
        }
        this.setState({
            imageIDs  : JSON.stringify(imageIDs),
            imageList : data.fileList
        });
    }
    onContentTypeChange(value) {
        this.setState({
            contentType: value
        });
    }
    onAddContent() {
        var contents = this.state.contents.slice(0);
        contents.push({
            id    : utils.uuid(),
            type  : this.state.contentType,
            value : ''
        });
        this.setState({
            contents: contents
        });
    }
    onContentTextChange(id, event) {
        var contents = this.state.contents.slice(0);
        for (var i = 0; i < contents.length; i++) {
            if (contents[i].id == id) {
                contents[i].value = event.target.value;
                this.setState({
                    contents: contents
                });
                break;
            }
        }
    }
    onContentImageChange(id, info) {
        if (info.file.status === 'done') {
            var contents = this.state.contents.slice(0);
            for (var i = 0; i < contents.length; i++) {
                if (contents[i].id == id) {
                    contents[i].value = info.file.response.data.url;
                    this.setState({
                        contents: contents
                    });
                    return;
                }
            }
        }
    }
    onSubmit() {
        const { dispatch } = this.props;
        dispatch(requestSaveProduct({
            id            : this.state.productId,
            name          : this.state.name,
            categories    : this.state.categories,
            status        : this.state.status,
            imageID       : this.state.imageID,
            imageIDs      : this.state.imageIDs,
            originalPrice : this.state.originalPrice,
            price         : this.state.price,
            remark        : this.state.remark,
            detail        : JSON.stringify(this.state.contents)
        }));
    }
    onPropValueVisibleChange(propId, visible) {
        var propValueVisibleMap = this.state.propValueVisibleMap;
        propValueVisibleMap[propId] = visible;
        this.setState({ 
            propValueVisibleMap: propValueVisibleMap,
            propValueTemp      : ''
        });
    }
    onPropValueInput(propId, event) {
        this.setState({ 
            propValueTemp: event.target.value
        });
    }
    addPropValue(propId) {
        var propValueTemp           = this.state.propValueTemp;
        var propValueVisibleMap     = this.state.propValueVisibleMap;
        propValueVisibleMap[propId] = false;
        this.setState({ 
            propValueVisibleMap : propValueVisibleMap,
            propValueTemp       : '' 
        });

        if (!propValueTemp) {
            return;
        }

        const { dispatch } = this.props;
        dispatch(requestSavePropertyValue({
            productID  : this.state.productId,
            propertyID : propId,
            name       : propValueTemp
        }));
    }
    cancelAddPropValue(propId) {
        var propValueVisibleMap     = this.state.propValueVisibleMap;
        propValueVisibleMap[propId] = false;
        this.setState({
            propValueVisibleMap : propValueVisibleMap,
            propValueTemp       : ''
        });
    }
    onPropVisibleChange(visible) {
        this.setState({ 
            propPopupVisible: visible
        });
    }
    onPropInput(event) {
        this.setState({
            propTemp: event.target.value
        });
    }
    addProp() {
        var propTemp = this.state.propTemp;
        this.setState({
            propTemp         : '',
            propPopupVisible : false
        });
        const { dispatch } = this.props;
        dispatch(requestSaveProperty({
            productID  : this.state.productId,
            name       : propTemp
        }));
    }
    cancelAddProp() {
        this.setState({
            propTemp         : '',
            propPopupVisible : false
        });
    }
    onInventoryChange(id, count) {
        const { dispatch } = this.props;
        dispatch(requestUpdateInventory({
            inventoryId : id,
            count       : count
        }));
    }
    onSaveInventory() {
        const { dispatch } = this.props;
        var inventories = [];
        for (var i = 0; i < this.state.inventories.length; i++) {
            inventories.push({
                id    : this.state.inventories[i].id,
                count : this.state.inventories[i].count
            });
        }
        dispatch(requestSaveInventory({
            productID   : this.state.productId,
            inventories : inventories
        }));
    }
    onTabChange(activeTabKey) {
        this.setState({
            activeTabKey: activeTabKey 
        });
    }
    onHasPropertyChange(value) {
        const { dispatch } = this.props;
        dispatch(requestUpdateHasProperty({
            productID : this.state.productId,
            value     : !!parseInt(value)
        }));
    }
    onTotalInventoryChange(count) {
        const { dispatch } = this.props;
        dispatch(requestUpdateTotalInventoryTemp({
            totalInventory : count
        }));
    }
    onSaveTotalInventory() {
        const { dispatch } = this.props;
        dispatch(requestUpdateTotalInventory({
            productID      : this.state.productId,
            totalInventory : this.state.totalInventory
        }));
    }
    render() {
        let self                = this;
        let { data }            = this.props;
        let editLabel           = this.state.productId ? '编辑' : '添加';
        let isLoading           = this.state.isLoading; 
        let name                = this.state.name;
        let contents            = this.state.contents;
        let originalPrice       = this.state.originalPrice;
        let price               = this.state.price;
        let remark              = this.state.remark;
        let status              = this.state.status;
        let imageData           = this.state.imageData;
        let uploadURL           = pageConfig.apiPath + '/admin/upload';
        let previewVisible      = this.state.previewVisible;
        let previewImage        = this.state.previewImage;
        let imageList           = this.state.imageList;
        let contentType         = this.state.contentType;
        let properties          = this.state.properties;
        let totalInventory      = this.state.totalInventory;
        let inventories         = this.state.inventories;
        let propValueVisibleMap = this.state.propValueVisibleMap;
        let propValueTemp       = this.state.propValueTemp;
        let propPopupVisible    = this.state.propPopupVisible;
        let propTemp            = this.state.propTemp;
        let activeTabKey        = this.state.activeTabKey;
        let hasProperty         = this.state.hasProperty ? 1 : 0;
        let hasPropertyValue    = this.state.hasPropertyValue;

        console.log(hasProperty, typeof hasProperty);

        let TabPane = Tabs.TabPane;

        let uploadButton = (
            <div>
                <Icon type="plus" />
                <div className="ant-upload-text">点击上传</div>
            </div>
        );

        const FormItem = Form.Item;
        const formItemLayout = {
            labelCol: {
                xs: { span: 24 },
                sm: { span: 3  },
            },
            wrapperCol: {
                xs: { span: 24 },
                sm: { span: 12 },
            }
        };

        const editorLayout = {
            labelCol: {
                xs: { span: 24 },
                sm: { span: 3  },
            },
            wrapperCol: {
                xs: { span: 24 },
                sm: { span: 21 },
            }
        };

        let treeData = utils.parseTree(data.categories);

        const treeProps = {
            treeData,
            value    : this.state.categories,
            onChange : this.onCategoriesChange,
            multiple : true,
            treeCheckable: true,
            showCheckedStrategy: TreeSelect.SHOW_PARENT,
            searchPlaceholder: '选择父分类',
            style: {
                width: '100%',
            }
        };

        return (
            <div>
                <Row gutter={24}>
                    <Col span={24}>
                        <div id="productBox">
                            <div className="product-title">{editLabel}商品</div>

                            <Tabs defaultActiveKey="1" onChange={self.onTabChange}>
                                <TabPane tab="商品信息" key="1">
                                {
                                    isLoading ? null :
                                    <Form>
                                        <FormItem {...formItemLayout} label="商品名称">
                                            <Input defaultValue={name} onBlur={this.onNameBlur}/>
                                        </FormItem>
                                        <FormItem {...formItemLayout} label="商品分类">
                                            <TreeSelect {...treeProps} />
                                        </FormItem>
                                        <FormItem {...formItemLayout} label="商品状态">
                                            <Select defaultValue={status} style={{ width: 120 }} onChange={this.onStatusChange}>
                                                <Select.Option value="3">等待上架</Select.Option>
                                                <Select.Option value="1">上架</Select.Option>
                                                <Select.Option value="2">下架</Select.Option>
                                            </Select>
                                        </FormItem>
                                        <FormItem {...formItemLayout} label="商品封面图">
                                            <Upload className="image-uploader" name="upFile"
                                                showUploadList={false} action={uploadURL}
                                                beforeUpload={this.onBeforeUpload}
                                                onChange={this.onImageChange}>
                                                {
                                                    imageData ?
                                                    <img src={imageData} alt="" className="image" /> 
                                                    :
                                                    <Icon type="plus" className="image-uploader-trigger" />
                                                }
                                            </Upload>
                                        </FormItem>
                                        <FormItem {...formItemLayout} label="商品图片集">
                                            <div className="clearfix">
                                                <Upload action={uploadURL} name="upFile"
                                                    listType="picture-card"
                                                    fileList={imageList}
                                                    onPreview={this.onPreview}
                                                    beforeUpload={this.onBeforeUpload}
                                                    onChange={this.onImageListChange}>
                                                { imageList.length >= 6 ? null : uploadButton }
                                                </Upload>
                                                <Modal visible={previewVisible} onCancel={this.onCancelPreview}
                                                    footer={null}>
                                                    <img alt="预览图片" style={{ width: '100%' }} src={previewImage} />
                                                </Modal>
                                            </div>
                                        </FormItem>
                                        <FormItem {...formItemLayout} label="原价">
                                            <InputNumber min={0} max={1000000} defaultValue={originalPrice} step={0.01} onBlur={this.onOriginalPriceBlur} />
                                            元
                                        </FormItem>
                                        <FormItem {...formItemLayout} label="促销价">
                                            <InputNumber min={0} max={1000000} defaultValue={price} step={0.01} onBlur={this.onPriceBlur} />
                                            元
                                        </FormItem>
                                        <FormItem {...formItemLayout} label="备注">
                                            <Input type="textarea" defaultValue={remark} rows={4} onBlur={this.onRemarkBlur}/>
                                        </FormItem>
                                        <FormItem {...editorLayout} label="商品详情">
                                        {
                                            contents.map(function(content) {
                                                return (
                                                    content.type == 'text' ?
                                                    <div key={content.id}>
                                                        <Input type="textarea" defaultValue={content.value} rows={4} onBlur={self.onContentTextChange.bind(self, content.id)}/>
                                                    </div>
                                                    :
                                                    <div key={content.id}>
                                                        <Upload className="image-uploader" name="upFile"
                                                            showUploadList={false} action={uploadURL}
                                                            beforeUpload={self.onBeforeUpload}
                                                            onChange={self.onContentImageChange.bind(self, content.id)}>
                                                            {
                                                                content.value ?
                                                                <img src={content.value} alt="" className="image" /> 
                                                                :
                                                                <Icon type="plus" className="image-uploader-trigger" />
                                                            }
                                                        </Upload>
                                                    </div>
                                                )
                                            })
                                        }
                                            <div>
                                                <Select defaultValue={contentType} style={{ width: 120 }} onChange={this.onContentTypeChange}>
                                                    <Select.Option value="image">图片</Select.Option>
                                                    <Select.Option value="text">文本</Select.Option>
                                                </Select>
                                                <Button onClick={this.onAddContent} type="primary" size="large">添加</Button>
                                            </div>
                                        </FormItem>
                                    </Form>
                                }
                                </TabPane>
                                <TabPane tab="商品库存" key="2">
                                    <Form>
                                        <FormItem className="produc" {...formItemLayout} label={'商品属性'}>
                                            <Select defaultValue={hasProperty + ''} style={{ width: 120 }} onChange={this.onHasPropertyChange}>
                                                <Select.Option value="0">无</Select.Option>
                                                <Select.Option value="1">有</Select.Option>
                                            </Select>
                                        </FormItem>
                                        {
                                            hasProperty ?
                                            properties.map(function(prop) {
                                                return (
                                                    <FormItem key={prop.id} {...formItemLayout} label={prop.name}>
                                                    {
                                                        prop.values.map(function(value) {
                                                            return (
                                                                <span key={value.id} className="product-prop-value">{value.name}</span>
                                                            )
                                                        })
                                                    }
                                                        <Popover content={
                                                            <div>
                                                                <Input value={propValueTemp} onChange={self.onPropValueInput.bind(self, prop.id)} className="product-prop-value-add-input"/>
                                                                <Button onClick={self.addPropValue.bind(self, prop.id)} type="primary" className="product-prop-value-add-confirm">确定</Button>
                                                                <Button onClick={self.cancelAddPropValue.bind(self, prop.id)}>取消</Button>
                                                            </div>} 
                                                            onVisibleChange={self.onPropValueVisibleChange.bind(self, prop.id)}
                                                            visible={propValueVisibleMap[prop.id]}
                                                            title={'添加' + prop.name} trigger="click" >
                                                            <Icon type="plus-circle" className="product-prop-value-add"/>
                                                        </Popover>
                                                    </FormItem>
                                                );
                                            })
                                            : ''
                                        }
                                        {
                                        hasProperty ?
                                        <FormItem className="product-prop-add" {...formItemLayout} label={' '}>
                                            <Popover content={
                                                <div>
                                                    <Input value={propTemp} onChange={self.onPropInput} className="product-prop-add-input"/>
                                                    <Button onClick={self.addProp} type="primary" className="product-prop-add-confirm">确定</Button>
                                                    <Button onClick={self.cancelAddProp}>取消</Button>
                                                </div>} 
                                                onVisibleChange={self.onPropVisibleChange}
                                                visible={propPopupVisible}
                                                title={'添加属性'} trigger="click" >
                                                <Button type="primary">添加属性</Button>
                                            </Popover>
                                        </FormItem>
                                        : ''
                                        }
                                        { 
                                        hasProperty && hasPropertyValue ? 
                                        <FormItem {...formItemLayout} label="库存">
                                        {
                                            inventories.map(function(inv) {
                                                var str = '';
                                                for (var i = 0; i < inv.propertyValues.length; i++) {
                                                    str += inv.propertyValues[i].name
                                                    if (i < inv.propertyValues.length - 1) {
                                                        str += '&nbsp;&nbsp;&nbsp;';
                                                    }
                                                }
                                                return (
                                                    <div key={inv.id} className="product-inventory-item">
                                                        <span className="product-inventory-label" dangerouslySetInnerHTML={{__html: str}}/>
                                                        <span className="product-inventory-unit">件</span>
                                                        <InputNumber onChange={self.onInventoryChange.bind(self, inv.id)} min={0} max={10000000000} defaultValue={inv.count} />
                                                    </div>
                                                )
                                            })
                                        }
                                            <div className="inventory-total">
                                                <span className="inventory-total-label">共</span>
                                                {totalInventory}
                                                <span className="inventory-total-unit">件</span>
                                            </div>
                                            <div className="inventory-save">
                                                <Button onClick={self.onSaveInventory} type="primary" size="large">保存库存</Button>
                                            </div>
                                        </FormItem>
                                        : ''
                                        }
                                        {
                                        !hasProperty ?
                                        <FormItem {...formItemLayout} label="总库存">
                                            <InputNumber className="inventory-total-input" min={0} max={10000000000} 
                                                onChange={self.onTotalInventoryChange}
                                                defaultValue={totalInventory} />
                                            <span>件</span>
                                        </FormItem>
                                        : ''
                                        }
                                        {
                                        !hasProperty ?
                                        <FormItem {...formItemLayout} label=" ">
                                            <Button onClick={self.onSaveTotalInventory} type="primary" size="large">保存库存</Button>
                                        </FormItem>
                                        : ''
                                        }
                                    </Form>
                                </TabPane>
                            </Tabs>
                        </div>
                    </Col>
                    <Col span={24} className="submit-box" style={{display: activeTabKey == '1' ? '' : 'none'}}>
                        <Button onClick={this.onSubmit} type="primary" size="large">保存</Button>
                        <Button className="submit-cancel-btn" size="large">取消</Button>
                    </Col>
                </Row>
                {/*<Software />*/}
            </div>
        );
    }
}

const mapStateToProps = (state) => {
    return {
        data: state.product
    };
}

export default connect(mapStateToProps)(EditProduct);

