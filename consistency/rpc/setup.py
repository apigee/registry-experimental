import io
import os
import setuptools  # type: ignore
from setuptools import find_packages

package_root = os.path.abspath(os.path.dirname(__file__))

setuptools.setup(
    name='google-cloud-apigeeregistry',
    packages=setuptools.PEP420PackageFinder.find(),
    install_requires=(
        'google-api-core[grpc] >= 1.31.0, < 3.0.0dev',
        'googleapis-common-protos >= 1.55.0, <2.0.0dev',
    ),
)
